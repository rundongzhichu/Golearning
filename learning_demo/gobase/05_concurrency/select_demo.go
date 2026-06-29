// Package concurrency 演示 Go 语言中的 select 语句。
//
// select 的特点：
//   - 用于等待多个 channel 操作
//   - 类似于 switch，但每个 case 必须是 channel 操作（发送或接收）
//   - 随机选择一个就绪的 case 执行（如果有多个同时就绪）
//   - 如果所有 case 都阻塞，则执行 default（如果有）
//   - 如果没有 default 且所有 case 都阻塞，select 会一直等待
package main

import (
	"fmt"
	"time"
)

// ==================== select 基础 ====================

func demoSelectBasics() {
	fmt.Println("=== select 基础 ===")

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	ch1 <- "hello"
	ch2 <- "world"

	// select 随机选择一个就绪的 case
	select {
	case msg := <-ch1:
		fmt.Printf("从 ch1: %s\n", msg)
	case msg := <-ch2:
		fmt.Printf("从 ch2: %s\n", msg)
	}
}

// ==================== select + timeout ====================

func demoSelectTimeout() {
	fmt.Println("\n=== select 超时模式 ===")

	ch := make(chan string)

	// 启动一个慢 goroutine
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- "slow response"
	}()

	select {
	case msg := <-ch:
		fmt.Printf("收到: %s\n", msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("超时！(100ms)")
	}
}

// ==================== select + default（非阻塞操作） ====================

func demoSelectDefault() {
	fmt.Println("\n=== select default（非阻塞） ===")

	ch := make(chan int, 1)

	// 非阻塞发送
	select {
	case ch <- 42:
		fmt.Println("发送成功")
	default:
		fmt.Println("channel 已满，跳过发送")
	}

	// 非阻塞接收
	select {
	case v := <-ch:
		fmt.Printf("接收成功: %d\n", v)
	default:
		fmt.Println("channel 为空，跳过接收")
	}
}

// ==================== select + done channel（取消） ====================

func demoSelectWithCancel() {
	fmt.Println("\n=== select 取消模式 ===")

	done := make(chan struct{})
	work := make(chan int)

	// 工作 goroutine
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("  收到取消信号，退出")
				return
			case n := <-work:
				fmt.Printf("  处理工作: %d\n", n)
			}
		}
	}()

	// 发送一些工作
	work <- 1
	work <- 2

	time.Sleep(50 * time.Millisecond)

	// 取消
	close(done)
	time.Sleep(50 * time.Millisecond)
}

// ==================== select + ticker（定时执行） ====================

func demoSelectWithTicker() {
	fmt.Println("\n=== select + Ticker 定时执行 ===")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan struct{})

	go func() {
		time.Sleep(500 * time.Millisecond)
		close(done)
	}()

	for {
		select {
		case <-done:
			fmt.Println("定时器结束")
			return
		case t := <-ticker.C:
			fmt.Printf("  tick at %v\n", t.Second())
		}
	}
}

// ==================== select 综合示例：心跳 + 超时 + 取消 ====================

func demoSelectComprehensive() {
	fmt.Println("\n=== select 综合示例 ===")

	// 模拟一个有超时、心跳和取消机制的任务
	task := func(
		done <-chan struct{},
		heartbeatInterval, timeout time.Duration,
	) <-chan string {
		result := make(chan string)

		go func() {
			defer close(result)

			heartbeat := time.NewTicker(heartbeatInterval)
			defer heartbeat.Stop()

			deadline := time.NewTimer(timeout)
			defer deadline.Stop()

			for {
				select {
				case <-done:
					fmt.Println("  [任务] 被外部取消")
					return
				case <-deadline.C:
					result <- "任务超时"
					return
				case <-heartbeat.C:
					fmt.Println("  💓 心跳...")
				default:
					// 模拟工作...
					time.Sleep(80 * time.Millisecond)
					result <- "任务完成"
					return
				}
			}
		}()

		return result
	}

	done := make(chan struct{})
	resultCh := task(done, 100*time.Millisecond, 2*time.Second)

	// 等待结果
	select {
	case result := <-resultCh:
		fmt.Printf("结果: %s\n", result)
	case <-time.After(3 * time.Second):
		fmt.Println("主程序超时")
	}
}

func main() {
	demoSelectBasics()
	demoSelectTimeout()
	demoSelectDefault()
	demoSelectWithCancel()
	demoSelectWithTicker()
	demoSelectComprehensive()
}
