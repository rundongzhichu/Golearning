// Package concurrency 演示 Go 语言中的 goroutine（轻量级协程）。
//
// goroutine 的特点：
//   - 使用 go 关键字启动，在独立协程中并发执行
//   - 比操作系统线程更轻量（初始栈约 2KB，动态伸缩）
//   - 由 Go 运行时调度（M:N 模型：M 个 goroutine 映射到 N 个 OS 线程）
//   - 通过 channel 或 sync 包进行通信和同步
//   - 推荐使用 CSP（Communicating Sequential Processes）模型
//   - GOMAXPROCS 控制同时执行的最大 CPU 数
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// ==================== 基础 goroutine ====================

// sayHello 在 goroutine 中执行的函数
func sayHello(name string) {
	fmt.Printf("Hello, %s! (goroutine)\n", name)
}

func demoBasicGoroutine() {
	fmt.Println("=== 基础 goroutine ===")

	// 启动一个 goroutine
	go sayHello("Alice")
	go sayHello("Bob")

	// 匿名函数 goroutine
	go func(msg string) {
		fmt.Printf("匿名 goroutine: %s\n", msg)
	}("Hello from closure")

	// 主 goroutine 等待一段时间，否则可能还没执行就退出了
	// 实际开发中应使用 sync.WaitGroup 或 channel 同步
	time.Sleep(100 * time.Millisecond)
}

// ==================== WaitGroup 同步 ====================

func demoWaitGroup1() {
	fmt.Println("\n=== WaitGroup 同步 ===")

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1) // 计数器 +1

		// 注意：必须传值给闭包，否则所有 goroutine 共享同一个 i
		go func(id int) {
			defer wg.Done() // goroutine 完成时计数器 -1
			time.Sleep(time.Duration(id*100) * time.Millisecond)
			fmt.Printf("  goroutine %d 完成\n", id)
		}(i)
	}

	wg.Wait() // 等待所有 goroutine 完成
	fmt.Println("所有 goroutine 已完成")
}

// ==================== GOMAXPROCS ====================

func demoGOMAXPROCS() {
	fmt.Println("\n=== GOMAXPROCS ===")

	fmt.Printf("CPU 核心数: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS:  %d\n", runtime.GOMAXPROCS(0))

	// 可以设置 GOMAXPROCS（通常不需要，让运行时自动管理）
	// runtime.GOMAXPROCS(2)
}

// ==================== goroutine 泄漏 ====================

// leakyGoroutine 演示 goroutine 泄漏：channel 阻塞导致 goroutine 永不退出
// 注意：这只是演示，实际使用中要避免
func demoGoroutineLeak() {
	fmt.Println("\n=== goroutine 泄漏演示 ===")

	// 正确的模式：使用 done channel 取消 goroutine
	done := make(chan struct{})

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				fmt.Println("  goroutine 收到取消信号，退出")
				return // 正确退出
			case t := <-ticker.C:
				fmt.Printf("  tick at %v\n", t.Second())
			}
		}
	}()

	time.Sleep(150 * time.Millisecond)
	close(done) // 通知 goroutine 退出
	time.Sleep(50 * time.Millisecond)
}

// ==================== goroutine 调度 ====================

// runtime.Gosched() 让出当前 goroutine 的执行权
func demoGosched() {
	fmt.Println("\n=== runtime.Gosched() ===")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			fmt.Printf("  goroutine A: %d\n", i)
			runtime.Gosched() // 让出 CPU 给其他 goroutine
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			fmt.Printf("  goroutine B: %d\n", i)
		}
	}()

	wg.Wait()
}

func main() {
	demoBasicGoroutine()
	demoWaitGroup()
	demoGOMAXPROCS()
	demoGoroutineLeak()
	demoGosched()
}
