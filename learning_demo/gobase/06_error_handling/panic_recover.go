// Package error_handling 演示 Go 语言中的 panic 和 recover。
//
// panic 和 recover 的设计哲学：
//   - 真正的异常情况才使用 panic（如数组越界、除零）
//   - 可预见的错误应使用 error 返回值处理
//   - recover 只在 defer 函数中有效
//   - 不应使用 panic/recover 来替代 try-catch
package main

import (
	"errors"
	"fmt"
)

// ==================== panic 触发 ====================

func causePanic() {
	fmt.Println("即将触发 panic...")
	panic("something went terribly wrong")
	// 下面的代码不会执行
	fmt.Println("这行永远不会执行")
}

// ==================== recover 捕获 ====================

func safeCall(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("⚠ recover 捕获到 panic: %v\n", r)
		}
	}()
	fn()
}

func demoRecover() {
	fmt.Println("=== recover 捕获 panic ===")

	fmt.Println("调用可能 panic 的函数（已保护）：")
	safeCall(causePanic)

	fmt.Println("\n函数返回后程序继续运行 ✓")
}

// ==================== panic + defer ====================

func demoPanicDefer() {
	fmt.Println("\n=== panic 与 defer ===")

	// defer 在 panic 时仍然会执行
	demo := func() {
		defer fmt.Println("  1. defer 语句 1 (总是执行)")
		defer fmt.Println("  2. defer 语句 2 (总是执行)")
		panic("突然 panic")
		// defer 语句 3 不会执行（因为它是函数末尾的匿名函数？不，defer 是在函数返回前执行）
	}

	safeCall(demo)
}

// ==================== recover 的返回值 ====================

func demoRecoverReturnValue() {
	fmt.Println("\n=== recover 的返回值 ===")

	// 以不同类型的值触发 panic，recover 可以获取到
	tests := []any{
		"string panic",
		42,
		errors.New("error panic"),
		struct{ Code int }{500},
	}

	for _, val := range tests {
		func() {
			defer func() {
				r := recover()
				fmt.Printf("  panic 值: %v (类型: %T)\n", r, r)
			}()
			panic(val)
		}()
	}
}

// ==================== 实际场景：Web 服务中间件 ====================

// recoverMiddleware 模拟 Web 框架的 panic 恢复中间件
func recoverMiddleware(handler func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  [中间件] 捕获 panic: %v, 返回 500 Internal Server Error\n", r)
		}
	}()
	handler()
}

func bootstrapService() {
	fmt.Println("\n=== Web 服务 panic 恢复模式 ===")

	// 正常请求
	fmt.Println("请求1: 正常处理")
	recoverMiddleware(func() {
		fmt.Println("  处理正常请求...")
	})

	// 会 panic 的请求
	fmt.Println("请求2: 处理失败")
	recoverMiddleware(func() {
		panic("数据库连接断开")
	})
}

// ==================== 禁止的用法 ====================

func demoAntiPatterns() {
	fmt.Println("\n=== 反模式（不推荐） ===")

	// 不要用 panic/recover 来模拟异常处理
	// 以下是不推荐的做法：
	badPattern := func() (result string) {
		defer func() {
			if r := recover(); r != nil {
				result = "从 panic 中恢复"
			}
		}()

		// 实际上应该返回 error，而不是 panic
		panic("业务错误")
	}

	fmt.Printf("  badPattern: %s\n", badPattern())
	fmt.Println("  提示: 业务错误应该返回 error，而不是 panic")
}

func main() {
	demoRecover()
	demoPanicDefer()
	demoRecoverReturnValue()
	bootstrapService()
	demoAntiPatterns()
}
