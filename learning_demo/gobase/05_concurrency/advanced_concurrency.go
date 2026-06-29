// Package concurrency 演示 Go 语言中的高级并发模式。
//
// 本文件涵盖：
//   - sync.WaitGroup：等待一组 goroutine 完成
//   - sync.Once：确保函数只执行一次
//   - sync.Pool：临时对象池（减少 GC 压力）
//   - sync.Cond：条件变量（不常用，但需要了解）
//   - sync/atomic：原子操作
//   - errgroup：错误组（golang.org/x/sync/errgroup）
//   - Semaphore 信号量模式
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ==================== sync.WaitGroup ====================

func demoWaitGroup() {
	fmt.Println("=== sync.WaitGroup ===")

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id*50) * time.Millisecond)
			fmt.Printf("  任务 %d 完成\n", id)
		}(i)
	}

	fmt.Println("等待所有任务完成...")
	wg.Wait()
	fmt.Println("全部完成")
}

// ==================== sync.Once ====================

type Database struct {
	connected bool
}

var (
	db     *Database
	once   sync.Once
	dbOnce sync.Once
)

// GetDatabase 单例模式：使用 sync.Once 确保只初始化一次
func GetDatabase() *Database {
	once.Do(func() {
		fmt.Println("  初始化数据库连接...")
		db = &Database{connected: true}
	})
	return db
}

func demoOnce() {
	fmt.Println("\n=== sync.Once ===")

	// 多次调用 GetDatabase，初始化函数只执行一次
	for i := 0; i < 3; i++ {
		db := GetDatabase()
		fmt.Printf("  第 %d 次获取: connected=%v\n", i+1, db.connected)
	}

	// Once 的妙用：延迟初始化、资源加载等
	var loadOnce sync.Once
	loadConfig := func() {
		loadOnce.Do(func() {
			fmt.Println("  加载配置文件...")
		})
	}

	loadConfig()
	loadConfig() // 第二次调用不会执行
}

// ==================== sync.Pool ====================

// BufferPool 字节缓冲区池
var bufferPool = sync.Pool{
	New: func() any {
		// 当池中没有可用对象时，创建新的
		return make([]byte, 1024)
	},
}

func demoPool() {
	fmt.Println("\n=== sync.Pool ===")

	// 获取对象
	buf := bufferPool.Get().([]byte)
	fmt.Printf("  从池获取 buffer: len=%d, cap=%d\n", len(buf), cap(buf))

	// 使用 buffer...
	copy(buf, []byte("hello"))

	// 放回池中
	bufferPool.Put(buf)

	// 再次获取（可能得到同一个对象）
	buf2 := bufferPool.Get().([]byte)
	fmt.Printf("  再次获取 buffer: len=%d, cap=%d\n", len(buf2), cap(buf2))

	// 注意：
	// 1. Pool 中的对象可能被 GC 自动清理
	// 2. 放入 Pool 前应该重置对象
	// 3. 适用于创建开销大、使用频繁的临时对象
}

// ==================== atomic 原子操作 ====================

func demoAtomic() {
	fmt.Println("\n=== sync/atomic 原子操作 ===")

	// 原子操作不需要锁，性能更高
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1) // 原子递增
		}()
	}

	wg.Wait()
	fmt.Printf("  原子计数器: %d\n", atomic.LoadInt64(&counter))

	// 原子比较并交换 (CAS)
	var flag int32
	swapped := atomic.CompareAndSwapInt32(&flag, 0, 1)
	fmt.Printf("  CAS: swapped=%v, flag=%d\n", swapped, flag)

	// 原子存储和读取
	atomic.StoreInt64(&counter, 100)
	fmt.Printf("  Store + Load: %d\n", atomic.LoadInt64(&counter))

	// Go 1.19+ 新增的原子类型
	var atomicVal atomic.Int64
	atomicVal.Store(42)
	fmt.Printf("  atomic.Int64: %d\n", atomicVal.Load())
	atomicVal.Add(10)
	fmt.Printf("  atomic.Int64 Add: %d\n", atomicVal.Load())
}

// ==================== 信号量模式 ====================

// Semaphore 简单信号量实现（使用缓冲 channel）
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore 创建信号量
func NewSemaphore(maxConcurrency int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, maxConcurrency),
	}
}

// Acquire 获取信号量
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release 释放信号量
func (s *Semaphore) Release() {
	<-s.ch
}

func demoSemaphore() {
	fmt.Println("\n=== 信号量模式 ===")

	// 限制同时最多 3 个 goroutine 执行
	sem := NewSemaphore(3)
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			fmt.Printf("  任务 %d 开始 (并发数: %d)\n", id, 3-len(sem.ch))
			time.Sleep(time.Duration(100) * time.Millisecond)
			fmt.Printf("  任务 %d 完成\n", id)
		}(i)
	}

	wg.Wait()
}

// ==================== sync.Cond ====================

func demoCond() {
	fmt.Println("\n=== sync.Cond 条件变量 ===")

	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	ready := false

	// 等待者
	go func() {
		mu.Lock()
		defer mu.Unlock()

		for !ready {
			fmt.Println("  等待者: 等待条件满足...")
			cond.Wait() // 释放锁并等待，被唤醒后重新获取锁
		}
		fmt.Println("  等待者: 条件满足，继续执行")
	}()

	// 触发者
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	ready = true
	mu.Unlock()
	cond.Signal() // 唤醒一个等待者（Broadcast 唤醒全部）

	time.Sleep(100 * time.Millisecond)
}

// ==================== 并发安全的集合 ====================

// SyncSet 线程安全的集合
type SyncSet struct {
	mu    sync.RWMutex
	items map[string]struct{}
}

// NewSyncSet 创建集合
func NewSyncSet() *SyncSet {
	return &SyncSet{items: make(map[string]struct{})}
}

// Add 添加元素
func (s *SyncSet) Add(item string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = struct{}{}
}

// Contains 检查是否包含
func (s *SyncSet) Contains(item string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.items[item]
	return ok
}

func demoSyncSet() {
	fmt.Println("\n=== 线程安全集合 ===")

	set := NewSyncSet()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			set.Add(fmt.Sprintf("item-%d", id%10))
		}(i)
	}

	wg.Wait()
	fmt.Printf("  集合大小: %d\n", len(set.items))
	fmt.Printf("  Contains 'item-5': %v\n", set.Contains("item-5"))
}

func main() {
	demoWaitGroup()
	demoOnce()
	demoPool()
	demoAtomic()
	demoSemaphore()
	demoCond()
	demoSyncSet()
}
