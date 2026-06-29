// Package methods_interfaces 演示 Go 语言中的接口（interface）。
//
// Go 接口的特点：
//   - 隐式实现：类型无需显式声明「实现了某接口」，只要方法集匹配即可
//   - 接口是抽象类型：只定义行为，不定义数据
//   - 接口值包含两部分：(类型, 值)
//   - 空接口 interface{} 可以存储任意类型的值（Go 1.18+ 推荐用 any）
//   - 最佳实践：接口要小（1-3 个方法），在使用方定义接口
//   - 零值为 nil：nil 接口的 (类型, 值) 都是 nil
package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

// ==================== 接口定义 ====================

// Shape 形状接口：定义面积计算
type Shape interface {
	Area() float64
}

// Stringer 打印接口（Go 标准库 fmt.Stringer）
type Stringer interface {
	String() string
}

// ==================== 接口实现 ====================

// Circle 圆
type Circle struct {
	Radius float64
}

// Area 实现 Shape 接口（值接收者）
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Rectangle 矩形
type Rectangle struct {
	Width, Height float64
}

// Area 实现 Shape 接口（值接收者）
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Triangle 三角形
type Triangle struct {
	Base, Height float64
}

// Area 实现 Shape 接口
func (t Triangle) Area() float64 {
	return t.Base * t.Height * 0.5
}

func demoInterfaceBasics() {
	fmt.Println("=== 接口基础 ===")

	// 接口变量可以存储任何实现了该接口的类型
	var s Shape

	s = Circle{Radius: 5}
	fmt.Printf("Circle Area: %.2f\n", s.Area())

	s = Rectangle{Width: 4, Height: 6}
	fmt.Printf("Rectangle Area: %.2f\n", s.Area())

	s = Triangle{Base: 3, Height: 4}
	fmt.Printf("Triangle Area: %.2f\n", s.Area())

	// 多态：统一处理不同的形状
	shapes := []Shape{
		Circle{Radius: 10},
		Rectangle{Width: 5, Height: 8},
		Triangle{Base: 6, Height: 4},
	}
	fmt.Println("\n多态遍历:")
	for i, shape := range shapes {
		fmt.Printf("  形状 %d: Area=%.2f\n", i, shape.Area())
	}
}

// ==================== 隐式接口实现 ====================

// 标准库的 io.Reader 接口：只需要实现 Read 方法
// type Reader interface { Read(p []byte) (n int, err error) }

// MyReader 自定义 Reader
type MyReader struct {
	data   string
	cursor int
}

// Read 实现 io.Reader 接口：隐式实现，无需声明 implements
func (r *MyReader) Read(p []byte) (n int, err error) {
	if r.cursor >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.cursor:])
	r.cursor += n
	return n, nil
}

func demoImplicitInterface() {
	fmt.Println("\n=== 隐式接口实现 ===")

	reader := &MyReader{data: "Hello, Go Interface!"}
	// MyReader 自动实现了 io.Reader，可以传递给任何需要 io.Reader 的地方
	buf := make([]byte, 8)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		fmt.Printf("  读取: %s\n", string(buf[:n]))
	}
}

// ==================== 接口组合 ====================

// ReadWriter 组合 Reader 和 Writer 接口
type ReadWriter interface {
	io.Reader
	io.Writer
}

// Closer 关闭接口
type Closer interface {
	Close() error
}

// ReadWriteCloser 组合多个接口
type ReadWriteCloser interface {
	io.Reader
	io.Writer
	Closer
}

func demoInterfaceComposition() {
	fmt.Println("\n=== 接口组合 ===")

	// strings.Builder 实现了 io.Writer 接口
	var w io.Writer = &strings.Builder{}
	fmt.Printf("strings.Builder 实现了 io.Writer: %T\n", w)

	// os.File 实现了 io.Reader、io.Writer、io.Closer 等多个接口
	var f *os.File // 实际中通过 os.Open 打开
	_ = f          // 演示类型关系

	var rwc ReadWriteCloser
	_ = rwc
}

// ==================== 接口值 ====================

func describe(i interface{}) {
	fmt.Printf("  类型: %T, 值: %v\n", i, i)
}

func demoInterfaceValues() {
	fmt.Println("\n=== 接口值 ===")

	var s Shape
	fmt.Printf("nil 接口: 类型=%T, 值=%v, isNil=%v\n", s, s, s == nil)

	// 接口赋值后：包含 (具体类型, 具体值)
	s = Circle{Radius: 3}
	fmt.Printf("赋值后: 类型=%T, 值=%v, isNil=%v\n", s, s, s == nil)

	// 空接口可以存储任何值
	var any interface{}
	any = 42
	describe(any)
	any = "hello"
	describe(any)
	any = []int{1, 2, 3}
	describe(any)
}

// ==================== 接口最佳实践 ====================

// 接口要小（单一职责）
type Reader interface {
	Read(p []byte) (n int, err error)
}

// 在使用方定义接口（而不是在实现方）
// consumer 消费者：只需定义自己需要的接口
type Greeter interface {
	Greet() string
}

// 实现方不知道 Greeter 接口的存在
type EnglishSpeaker struct{}

func (EnglishSpeaker) Greet() string { return "Hello!" }

type ChineseSpeaker struct{}

func (ChineseSpeaker) Greet() string { return "你好!" }

func greetEveryone(g Greeter) {
	fmt.Println(g.Greet())
}

func demoInterfaceBestPractices() {
	fmt.Println("\n=== 接口最佳实践 ===")

	// 小接口组合而不是大接口
	// 在使用方定义接口
	greetEveryone(EnglishSpeaker{})
	greetEveryone(ChineseSpeaker{})
}

func main() {
	demoInterfaceBasics()
	demoImplicitInterface()
	demoInterfaceComposition()
	demoInterfaceValues()
	demoInterfaceBestPractices()
}
