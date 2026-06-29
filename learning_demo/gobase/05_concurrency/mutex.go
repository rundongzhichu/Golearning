// Package concurrency 演示 Go 语言中的互斥锁（Mutex 和 RWMutex）。
//
// 锁的特点：
//   - sync.Mutex：互斥锁，同一时刻只有一个 goroutine 能持有锁
//   - sync.RWMutex：读写锁，多个读操作可并发，写操作独占
//   - 推荐与 defer 配合使用，确保解锁
//   - 锁的零值就是可用的（不需要初始化）
//   - 锁不应该被复制（值传递）
package main

import (
	"fmt"
	"sync"
)

// ==================== Mutex 基础 ====================

// SafeCounter 线程安全的计数器
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

// Increment 增加计数（线程安全）
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// Value 获取当前值（线程安全）
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func demoMutexBasics() {
	fmt.Println("=== Mutex 基础 ===")

	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// 100 个 goroutine 并发递增
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("最终计数: %d (预期 100)\n", counter.Value())
}

// ==================== 不加锁的竞态条件 ====================

func demoRaceCondition() {
	fmt.Println("\n=== 竞态条件演示 ===")

	var count int // 不加锁的计数器
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++ // 竞态！多个 goroutine 同时读写
		}()
	}

	wg.Wait()
	fmt.Printf("最终计数: %d (预期 1000，实际可能小于1000)\n", count)
	fmt.Println("提示: 使用 go run -race 检测竞态条件")
}

// ==================== RWMutex 读写锁 ====================

// SafeCache 线程安全的缓存（读写锁）
type SafeCache struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewSafeCache 创建缓存
func NewSafeCache() *SafeCache {
	return &SafeCache{
		store: make(map[string]string),
	}
}

// Get 读取（使用读锁，允许多个并发读取）
func (c *SafeCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.store[key]
	return value, ok
}

// Set 写入（使用写锁，独占）
func (c *SafeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

// Delete 删除（使用写锁）
func (c *SafeCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func demoRWMutex() {
	fmt.Println("\n=== RWMutex 读写锁 ===")

	cache := NewSafeCache()
	var wg sync.WaitGroup

	// 写入
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Set(key, fmt.Sprintf("value-%d", i))
	}

	// 并发读取（读锁不互斥）
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key-%d", j%5)
				if value, ok := cache.Get(key); ok {
					_ = value
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("并发读取完成")

	// 测试读取
	for i := 0; i < 5; i++ {
		if v, ok := cache.Get(fmt.Sprintf("key-%d", i)); ok {
			fmt.Printf("  %s: %s\n", fmt.Sprintf("key-%d", i), v)
		}
	}
}

// ==================== TryLock（Go 1.18+） ====================

func demoTryLock() {
	fmt.Println("\n=== TryLock (Go 1.18+) ===")

	var mu sync.Mutex

	// 尝试加锁，不阻塞
	if mu.TryLock() {
		fmt.Println("第一次 TryLock 成功")
		mu.Unlock()
	}

	mu.Lock()
	// 已经锁住了，TryLock 会失败
	if !mu.TryLock() {
		fmt.Println("第二次 TryLock 失败（已被锁住）")
	}
	mu.Unlock()
}

// ==================== sync.Map ====================

func demoSyncMap() {
	fmt.Println("\n=== sync.Map ===")

	var sm sync.Map

	// 写入
	sm.Store("key1", "value1")
	sm.Store("key2", "value2")

	// 读取
	if v, ok := sm.Load("key1"); ok {
		fmt.Printf("Load key1: %v\n", v)
	}

	// 读取或写入
	v, loaded := sm.LoadOrStore("key3", "value3")
	fmt.Printf("LoadOrStore key3: %v, loaded=%v\n", v, loaded)

	// 遍历
	sm.Range(func(key, value any) bool {
		fmt.Printf("  %v: %v\n", key, value)
		return true // 返回 false 停止遍历
	})

	// 删除
	sm.Delete("key1")
}

// ==================== 死锁演示 ====================

func demoDeadLock() {
	fmt.Println("\n=== 死锁演示 ===")

	type Account struct {
		sync.Mutex
		Balance int
	}

	// 转账函数（正确实现：按固定顺序加锁避免死锁）
	transfer := func(from, to *Account, amount int) {
		// 始终按地址从小到大加锁，避免死锁
		if fmt.Sprintf("%p", from) < fmt.Sprintf("%p", to) {
			from.Lock()
			to.Lock()
		} else {
			to.Lock()
			from.Lock()
		}
		defer from.Unlock()
		defer to.Unlock()

		from.Balance -= amount
		to.Balance += amount
	}

	a := &Account{Balance: 1000}
	b := &Account{Balance: 1000}

	fmt.Printf("转账前: a=%d, b=%d\n", a.Balance, b.Balance)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transfer(a, b, 10)
		}()
	}

	wg.Wait()
	fmt.Printf("转账后: a=%d, b=%d (总额=%d)\n", a.Balance, b.Balance, a.Balance+b.Balance)
}

func main() {
	demoMutexBasics()
	demoRaceCondition()
	demoRWMutex()
	demoTryLock()
	demoSyncMap()
	demoDeadLock()
}
