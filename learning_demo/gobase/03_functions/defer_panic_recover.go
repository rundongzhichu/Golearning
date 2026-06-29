// Package functions 演示 Go 语言中的 defer、panic 和 recover。
//
// defer：延迟执行
//   - 在函数返回前执行，常用于资源清理（关闭文件、释放锁等）
//   - 多个 defer 按 LIFO（后进先出）顺序执行
//   - defer 的参数在 defer 语句执行时求值（不是最终调用时）
//
// panic：运行时恐慌
//   - 类似于其他语言的异常
//   - 会立即停止当前函数的执行，并向上传递
//   - 沿调用栈向上，直到被 recover 捕获或程序崩溃
//
// recover：恢复恐慌
//   - 只能在 defer 函数中使用
//   - 捕获 panic 的值，防止程序崩溃
//   - 返回 panic 的值，若没有 panic 则返回 nil
package main

import (
	"fmt"
	"os"
	"sync"
)

// ==================== defer 基础 ====================

func demoDeferBasics() {
	fmt.Println("=== defer 基础 ===")

	// defer 在函数返回前执行
	defer fmt.Println("defer: 这句最后执行")
	fmt.Println("正常: 这句先执行")

	// 资源清理模式：打开文件后 defer 关闭
	readFile := func(filename string) error {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close() // 确保文件在函数返回时关闭

		// 读取文件...
		fmt.Println("  文件已打开，即将关闭...")
		return nil
	}

	_ = readFile("/tmp/test.txt") // 忽略错误
}

// ==================== defer 执行顺序 ====================

func demoDeferOrder() {
	fmt.Println("\n=== defer 顺序 (LIFO) ===")

	// defer 按后进先出（LIFO）顺序执行
	defer fmt.Println("defer 1")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3")
	fmt.Println("这是普通语句")

	// defer 的参数在 defer 语句执行时求值
	x := 1
	defer fmt.Printf("defer x=%d (defer语句时求值: 1)\n", x)
	x = 2
	fmt.Printf("赋值后 x=%d\n", x)
	// 输出：defer x=1（不是2！）
}

// ==================== 互斥锁 defer 模式 ====================

func demoDeferWithMutex() {
	fmt.Println("\n=== defer 与互斥锁 ===")

	var mu sync.Mutex
	counter := 0

	increment := func() {
		mu.Lock()
		defer mu.Unlock() // 确保解锁
		counter++
		fmt.Printf("  counter=%d\n", counter)
	}

	increment()
	increment()
}

// ==================== panic 和 recover ====================

// safeDivide 安全的除法：捕获 panic
func safeDivide(a, b int) (result int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  ⚠ recover 捕获: %v\n", r)
			result = 0 // 设置默认返回值
		}
	}()

	if b == 0 {
		panic("除数为零！") // 触发 panic
	}
	return a / b
}

// nestedPanic 演示 panic 的传播
func nestedPanic() {
	panic("内部 panic：something went wrong")
}

func outerFunction() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  外部函数 recover: %v\n", r)
		}
	}()
	fmt.Println("  调用 nestedPanic 前")
	nestedPanic()
	fmt.Println("  这行不会执行")
}

func demoPanicRecover() {
	fmt.Println("\n=== panic 和 recover ===")

	fmt.Println("safeDivide(10, 2) =", safeDivide(10, 2))
	fmt.Println("safeDivide(10, 0) =", safeDivide(10, 0))

	fmt.Println("\n  panic 传播示例：")
	outerFunction()
	fmt.Println("  outerFunction 返回后继续执行")

	// 注意：recover 只能在 defer 中使用
	// func() { recover() }()  // 无效！
}

// ==================== 自定义错误处理 ====================

// MyError 自定义错误类型
type MyError struct {
	Op  string
	Err error
}

func (e *MyError) Error() string {
	return fmt.Sprintf("操作 %s 失败: %v", e.Op, e.Err)
}

func doSomething() error {
	return &MyError{Op: "数据库查询", Err: fmt.Errorf("连接超时")}
}

func demoErrorHandling() {
	fmt.Println("\n=== 错误处理模式 ===")

	// 模式1：if err != nil
	result, err := os.ReadFile("/nonexistent")
	if err != nil {
		fmt.Printf("读取文件错误: %v\n", err)
	} else {
		fmt.Printf("内容: %s\n", result)
	}

	// 模式2：errors.Is 和 errors.As（另见 error_handling 目录）
	err = doSomething()
	fmt.Printf("自定义错误: %v\n", err)
}

func main() {
	demoDeferBasics()
	demoDeferOrder()
	demoDeferWithMutex()
	demoPanicRecover()
	demoErrorHandling()
}
