// Package concurrency 演示 Go 语言中的 channel（通道）。
//
// channel 的特点：
//   - Go 并发编程的核心：用于 goroutine 之间的通信
//   - 类型安全：chan T 只能发送/接收 T 类型
//   - 可以使用 make(chan T) 创建无缓冲 channel
//   - 使用 make(chan T, n) 创建缓冲 channel（容量 n）
//   - 发送：ch <- value；接收：value := <-ch
//   - 关闭：close(ch)，接收方可以检测到关闭
//   - 设计原则：Don't communicate by sharing memory; share memory by communicating.
//   - 单向 channel：chan<- T（只发送），<-chan T（只接收）
package main

import (
	"fmt"
	"sync"
	"time"
)

// ==================== 无缓冲 channel ====================

func demoUnbufferedChannel() {
	fmt.Println("=== 无缓冲 channel ===")

	// 无缓冲 channel：发送方必须等接收方准备好（同步）
	ch := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "ping" // 阻塞，直到接收方准备好
		fmt.Println("  发送完成")
	}()

	msg := <-ch // 阻塞，直到收到数据
	fmt.Printf("  接收到: %s\n", msg)
}

// ==================== 缓冲 channel ====================

func demoBufferedChannel() {
	fmt.Println("\n=== 缓冲 channel ===")

	// 缓冲 channel：缓冲区满之前不阻塞
	ch := make(chan int, 3)

	// 发送（缓冲区未满，不阻塞）
	ch <- 1
	ch <- 2
	ch <- 3
	// ch <- 4 // 缓冲区满了，会阻塞

	fmt.Printf("len=%d, cap=%d\n", len(ch), cap(ch))

	// 接收
	fmt.Printf("<-ch: %d\n", <-ch)
	fmt.Printf("<-ch: %d\n", <-ch)
	fmt.Printf("<-ch: %d\n", <-ch)

	fmt.Printf("接收完后: len=%d, cap=%d\n", len(ch), cap(ch))
}

// ==================== channel 关闭与遍历 ====================

func demoChannelClose() {
	fmt.Println("\n=== channel 关闭 ===")

	ch := make(chan int, 5)

	// 生产数据
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch) // 关闭 channel
		fmt.Println("  channel 已关闭")
	}()

	// 方式1：for range 遍历（自动检测关闭）
	fmt.Print("  for range: ")
	for v := range ch {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// 方式2：comma-ok 判断是否关闭
	ch2 := make(chan int, 1)
	ch2 <- 100
	close(ch2)

	v, ok := <-ch2
	fmt.Printf("  从已关闭 channel 接收: v=%d, ok=%v\n", v, ok)

	// 从已关闭 channel 再次接收：返回零值
	v, ok = <-ch2
	fmt.Printf("  再次接收: v=%d, ok=%v\n", v, ok)
}

// ==================== 单向 channel ====================

// producer 只发送 channel
func producer(ch chan<- int) {
	for i := 1; i <= 3; i++ {
		ch <- i
	}
	close(ch)
}

// consumer 只接收 channel
func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Printf("  消费: %d\n", v)
	}
}

func demoDirectionalChannel() {
	fmt.Println("\n=== 单向 channel ===")

	ch := make(chan int, 3)

	go producer(ch) // chan int 自动转为 chan<- int
	consumer(ch)    // chan int 自动转为 <-chan int
}

// ==================== select 多路复用 ====================

func demoSelectChannel() {
	fmt.Println("\n=== channel + select ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// 两个生产者
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "来自 ch1"
	}()
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "来自 ch2"
	}()

	// select 等待多个 channel
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("  收到: %s\n", msg)
		case msg := <-ch2:
			fmt.Printf("  收到: %s\n", msg)
		}
	}
}

// ==================== channel 常见模式 ====================

// 模式1：生产者-消费者
func demoProducerConsumer() {
	fmt.Println("\n=== 生产者-消费者模式 ===")

	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// 3 个消费者（worker）
	for w := 1; w <= 3; w++ {
		go func(id int) {
			for job := range jobs {
				fmt.Printf("  worker %d 处理 job %d\n", id, job)
				time.Sleep(50 * time.Millisecond)
				results <- job * 2
			}
		}(w)
	}

	// 生产者：发送 5 个任务
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// 收集结果
	for i := 1; i <= 5; i++ {
		fmt.Printf("  结果: %d\n", <-results)
	}
}

// 模式2：扇出-扇入（Fan-out/Fan-in）
func demoFanOutIn() {
	fmt.Println("\n=== 扇出-扇入模式 ===")

	// 输入 channel
	input := make(chan int, 10)

	// 扇出：多个 goroutine 从同一个 channel 消费
	outputs := make([]chan int, 3)
	for i := 0; i < 3; i++ {
		outputs[i] = make(chan int, 10)
		go func(id int, out chan int) {
			for n := range input {
				out <- n * n // 每个 worker 做平方运算
				time.Sleep(20 * time.Millisecond)
			}
			close(out)
		}(i, outputs[i])
	}

	// 发送数据
	go func() {
		for i := 1; i <= 6; i++ {
			input <- i
		}
		close(input)
	}()

	// 扇入：合并多个 channel 到一个 channel
	merged := mergeChannels(outputs...)

	for v := range merged {
		fmt.Printf("  扇入结果: %d\n", v)
	}
}

// mergeChannels 合并多个 channel
func mergeChannels(channels ...chan int) chan int {
	out := make(chan int)
	var wg sync.WaitGroup // 注意：这里实际是 goroutine 级别

	for _, ch := range channels {
		wg.Add(1)
		go func(c chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	demoUnbufferedChannel()
	demoBufferedChannel()
	demoChannelClose()
	demoDirectionalChannel()
	demoSelectChannel()
	demoProducerConsumer()
	demoFanOutIn()
}
