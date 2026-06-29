// Package concurrency 演示 Go 语言中的 context（上下文）。
//
// context 的特点：
//   - 用于在 goroutine 之间传递请求范围内的值、取消信号和截止时间
//   - 树形结构：父 context 取消时，所有子 context 也被取消
//   - context.Context 是一个接口，核心方法：
//   - Deadline()：返回截止时间
//   - Done()：返回一个只读 channel，context 取消时会关闭
//   - Err()：context 取消的原因（Canceled 或 DeadlineExceeded）
//   - Value(key)：返回 context 中存储的值
//   - context 是 Go 中处理超时、取消、传递元数据的标准方式
//   - 四个创建函数：
//   - context.Background()：根 context（最常用）
//   - context.TODO()：不确定用哪个 context 时的占位符
//   - context.WithCancel(parent)：可取消
//   - context.WithDeadline(parent, time)：有截止时间
//   - context.WithTimeout(parent, duration)：有超时时间
//   - context.WithValue(parent, key, value)：携带键值对
package main

import (
	"context"
	"fmt"
	"time"
)

// ==================== context.WithCancel ====================

// worker 工作函数：监听 context 取消信号
func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  worker %d: 收到取消信号，退出。原因: %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("  worker %d: 工作中...\n", id)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func demoWithCancel() {
	fmt.Println("=== context.WithCancel ===")

	ctx, cancel := context.WithCancel(context.Background())

	// 启动 3 个 worker
	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	// 工作 500ms 后取消
	time.Sleep(500 * time.Millisecond)
	fmt.Println("主 goroutine: 发送取消信号")
	cancel()

	// 等待 worker 退出
	time.Sleep(200 * time.Millisecond)
}

// ==================== context.WithTimeout ====================

func demoWithTimeout() {
	fmt.Println("\n=== context.WithTimeout ===")

	// 设置 300ms 超时
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel() // 确保释放资源

	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("任务完成")
	case <-ctx.Done():
		fmt.Printf("超时: %v\n", ctx.Err())
	}
}

// ==================== context.WithDeadline ====================

func demoWithDeadline() {
	fmt.Println("\n=== context.WithDeadline ===")

	// 设置绝对截止时间：当前时间 + 300ms
	deadline := time.Now().Add(300 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("任务完成")
	case <-ctx.Done():
		fmt.Printf("超过截止时间: %v\n", ctx.Err())
	}
}

// ==================== context.WithValue ====================

// contextKey 定义私有类型作为 context key（最佳实践：避免 key 冲突）
type contextKey string

const (
	KeyUserID    contextKey = "userID"
	KeyRequestID contextKey = "requestID"
	KeyTraceID   contextKey = "traceID"
)

// processRequest 处理请求：从 context 中读取值
func processRequest(ctx context.Context) {
	userID, ok := ctx.Value(KeyUserID).(string)
	if !ok {
		userID = "unknown"
	}
	reqID, _ := ctx.Value(KeyRequestID).(string)
	traceID, _ := ctx.Value(KeyTraceID).(string)

	fmt.Printf("  处理请求: userID=%s, requestID=%s, traceID=%s\n",
		userID, reqID, traceID)
}

func demoWithValue() {
	fmt.Println("\n=== context.WithValue ===")

	// 构建带值的 context 链
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyUserID, "user-123")
	ctx = context.WithValue(ctx, KeyRequestID, "req-456")
	ctx = context.WithValue(ctx, KeyTraceID, "trace-789")

	processRequest(ctx)
}

// ==================== context 取消传播 ====================

func demoCancelPropagation() {
	fmt.Println("\n=== context 取消传播 ===")

	// 父 context 取消时，所有子 context 都取消
	parentCtx, cancel := context.WithCancel(context.Background())

	// 创建子 context
	childCtx, _ := context.WithCancel(parentCtx)

	go func() {
		<-childCtx.Done()
		fmt.Printf("  子 context: %v\n", childCtx.Err())
	}()

	time.Sleep(100 * time.Millisecond)
	cancel() // 取消父 context，子 context 也会取消

	time.Sleep(50 * time.Millisecond)
}

// ==================== HTTP 请求超时模式 ====================

// simulateHTTPRequest 模拟 HTTP 请求（context 超时模式）
func simulateHTTPRequest(ctx context.Context, url string, duration time.Duration) error {
	resultCh := make(chan string, 1)

	go func() {
		// 模拟请求处理
		time.Sleep(duration)
		resultCh <- fmt.Sprintf("响应来自 %s (耗时 %v)", url, duration)
	}()

	select {
	case result := <-resultCh:
		fmt.Printf("  成功: %s\n", result)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("请求 %s 被取消: %v", url, ctx.Err())
	}
}

func demoHTTPPattern() {
	fmt.Println("\n=== HTTP 请求超时模式 ===")

	ctx := context.Background()

	// 正常请求（100ms，在超时内完成）
	ctx1, cancel1 := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel1()
	err := simulateHTTPRequest(ctx1, "/api/fast", 100*time.Millisecond)
	if err != nil {
		fmt.Printf("  err: %v\n", err)
	}

	// 慢请求（500ms，超时）
	ctx2, cancel2 := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel2()
	err = simulateHTTPRequest(ctx2, "/api/slow", 500*time.Millisecond)
	if err != nil {
		fmt.Printf("  err: %v\n", err)
	}
}

func main() {
	demoWithCancel()
	demoWithTimeout()
	demoWithDeadline()
	demoWithValue()
	demoCancelPropagation()
	demoHTTPPattern()
}
