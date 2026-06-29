// Package functions 演示 Go 语言中的函数定义与使用。
//
// Go 函数的特点：
//   - 函数是一等公民：可以作为参数、返回值、赋值给变量
//   - 支持多返回值
//   - 支持命名返回值
//   - 支持可变参数（variadic）
//   - 支持匿名函数和闭包
//   - 函数本身没有默认参数和函数重载
package main

import (
	"errors"
	"fmt"
)

// ==================== 函数定义 ====================

// add 基础函数：两个 int 参数，返回 int
func add(a, b int) int {
	return a + b
}

// swap 多返回值
func swap(a, b string) (string, string) {
	return b, a
}

// divide 多返回值 + 错误处理
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}

// split 命名返回值：返回值名称直接作为函数内的局部变量
func split(sum int) (x, y int) {
	x = sum * 4 / 9 // 计算 x
	y = sum - x     // 计算 y
	return          // 裸返回（naked return）：自动返回 x 和 y
}

// sum 可变参数（variadic）：接收任意数量的 int
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func demoFunctionBasics() {
	fmt.Println("=== 函数基础 ===")

	fmt.Printf("add(3, 5) = %d\n", add(3, 5))

	x, y := swap("hello", "world")
	fmt.Printf("swap: %s, %s\n", x, y)

	result, err := divide(10, 3)
	fmt.Printf("divide(10,3) = %.2f, err=%v\n", result, err)
	result, err = divide(10, 0)
	fmt.Printf("divide(10,0) = %.2f, err=%v\n", result, err)

	a, b := split(17)
	fmt.Printf("split(17): x=%d, y=%d\n", a, b)

	fmt.Printf("sum(1,2,3,4,5) = %d\n", sum(1, 2, 3, 4, 5))
	// 展开切片传递给可变参数
	nums := []int{1, 2, 3}
	fmt.Printf("sum(nums...) = %d\n", sum(nums...))
}

// ==================== 函数类型 ====================

// Operation 定义函数类型
type Operation func(a, b int) int

// calculate 接收函数类型作为参数
func calculate(a, b int, op Operation) int {
	return op(a, b)
}

// getOperation 返回函数类型（闭包）
func getOperation(op string) Operation {
	switch op {
	case "add":
		return func(a, b int) int { return a + b }
	case "sub":
		return func(a, b int) int { return a - b }
	default:
		return nil
	}
}

func demoFunctionType() {
	fmt.Println("\n=== 函数类型 ===")

	// 函数赋值给变量
	multiply := func(a, b int) int {
		return a * b
	}
	fmt.Printf("multiply(3, 4) = %d\n", multiply(3, 4))

	// 匿名函数直接调用
	result := func(x int) int {
		return x * x
	}(5)
	fmt.Printf("匿名函数直接调用: 5² = %d\n", result)

	// 函数作为参数
	fmt.Printf("calculate(10,5, add) = %d\n", calculate(10, 5, add))

	// 函数作为返回值
	addOp := getOperation("add")
	fmt.Printf("getOperation('add')(10,5) = %d\n", addOp(10, 5))
}

// ==================== 闭包 ====================

// counter 返回一个计数器闭包
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// makeAdder 返回一个加法器闭包（捕获外部参数）
func makeAdder(n int) func(int) int {
	return func(x int) int {
		return x + n
	}
}

// fibonacci 返回一个 Fibonacci 数列生成器
func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func demoClosures() {
	fmt.Println("\n=== 闭包 ===")

	// 计数器闭包
	cnt := counter()
	fmt.Printf("counter: %d, %d, %d\n", cnt(), cnt(), cnt())

	// 每个闭包有独立的状态
	cnt2 := counter()
	fmt.Printf("另一个 counter: %d, %d\n", cnt2(), cnt2())
	fmt.Printf("原 counter 继续: %d\n", cnt()) // 状态独立

	// 加法器闭包
	add5 := makeAdder(5)
	add10 := makeAdder(10)
	fmt.Printf("add5(3)=%d, add10(3)=%d\n", add5(3), add10(3))

	// Fibonacci 生成器
	fib := fibonacci()
	fmt.Print("Fibonacci: ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", fib())
	}
	fmt.Println()
}

// ==================== init 函数 ====================

// init 在 main 之前自动执行，每个包可以有多个 init
// 执行顺序：导入的包的 init → 当前包的 init → main
func init() {
	// 常用于初始化配置、注册驱动等
	fmt.Println("[init] 包初始化完成")
}

func main() {
	demoFunctionBasics()
	demoFunctionType()
	demoClosures()
}
