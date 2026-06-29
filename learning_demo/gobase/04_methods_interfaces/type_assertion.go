// Package methods_interfaces 演示 Go 语言中的类型断言与类型 switch。
//
// 类型断言（type assertion）：
//   - 用于将接口类型转换为具体类型
//   - 语法：x.(T)，x 是接口值，T 是具体类型
//   - comma-ok 语法：v, ok := x.(T)，安全断言不会 panic
//
// 类型 switch：
//   - 根据接口值的动态类型执行不同的分支
//   - 语法：switch v := x.(type)
package main

import (
	"fmt"
)

// ==================== 类型断言 ====================

func demoTypeAssertion() {
	fmt.Println("=== 类型断言 ===")

	var i interface{} = "hello"

	// 不安全断言：如果类型不匹配会 panic
	s := i.(string)
	fmt.Printf("断言为 string: %s\n", s)

	// 安全断言：comma-ok 语法
	f, ok := i.(float64)
	fmt.Printf("断言为 float64: 值=%v, ok=%v\n", f, ok)

	// 错误类型断言会导致 panic（仅在确定类型时使用不安全断言）
	// i2 := i.(int) // panic: interface conversion: interface {} is string, not int
}

// ==================== 类型开关 ====================

// Animal 动物接口
type Animal interface {
	Speak() string
}

type Dog struct{ Name string }

func (d Dog) Speak() string { return d.Name + ": 汪汪!" }

type Cat struct{ Name string }

func (c Cat) Speak() string { return c.Name + ": 喵喵!" }

type Bird struct{ Name string }

func (b Bird) Speak() string { return b.Name + ": 啾啾!" }

// describeAnimal 类型 switch：根据具体类型做不同处理
func describeAnimal(a Animal) {
	fmt.Printf("描述 %s: ", a.Speak())

	// 类型 switch 语法
	switch v := a.(type) {
	case Dog:
		fmt.Printf("这是一只狗，名字叫 %s\n", v.Name)
	case Cat:
		fmt.Printf("这是一只猫，名字叫 %s\n", v.Name)
	case Bird:
		fmt.Printf("这是一只鸟，名字叫 %s\n", v.Name)
	default:
		fmt.Printf("未知动物\n")
	}
}

func demoTypeSwitch() {
	fmt.Println("\n=== 类型 switch ===")

	animals := []Animal{
		Dog{Name: "旺财"},
		Cat{Name: "咪咪"},
		Bird{Name: "小飞"},
	}

	for _, a := range animals {
		describeAnimal(a)
	}
}

// ==================== 高级类型 switch ====================

// doSomething 演示对空接口的类型 switch
func doSomething(x interface{}) {
	switch v := x.(type) {
	case int:
		fmt.Printf("int: %d (两倍: %d)\n", v, v*2)
	case string:
		fmt.Printf("string: %s (长度: %d)\n", v, len(v))
	case bool:
		if v {
			fmt.Println("bool: true")
		} else {
			fmt.Println("bool: false")
		}
	case []int:
		fmt.Printf("[]int: %v (len=%d)\n", v, len(v))
	case nil:
		fmt.Println("nil")
	case fmt.Stringer:
		fmt.Printf("fmt.Stringer: %s\n", v.String())
	default:
		fmt.Printf("未知类型: %T, 值: %v\n", v, v)
	}
}

func demoAdvancedTypeSwitch() {
	fmt.Println("\n=== 高级类型 switch ===")

	doSomething(42)
	doSomething("hello")
	doSomething(true)
	doSomething([]int{1, 2, 3})
	doSomething(nil)
}

// ==================== 常见模式 ====================

// 模式1：判断接口是否实现了某接口
type Writer interface {
	Write(p []byte) (n int, err error)
}

// 模式2：提取错误类型
type Temporary interface {
	Temporary() bool
}

func isTemporary(err error) bool {
	// 类型断言判断错误是否实现了 Temporary 接口
	type temp interface {
		Temporary() bool
	}
	t, ok := err.(temp)
	return ok && t.Temporary()
}

func demoPatterns() {
	fmt.Println("\n=== 常见模式 ===")

	// 使用类型断言判断是否实现了某个接口
	var v interface{} = "hello"
	if s, ok := v.(fmt.Stringer); ok {
		fmt.Printf("实现了 Stringer: %s\n", s.String())
	} else {
		fmt.Println("未实现 Stringer")
	}

	// 确保编译时类型实现了某接口（惯用法）
	// var _ io.Reader = (*MyType)(nil)  // 编译期检查
}

func main() {
	demoTypeAssertion()
	demoTypeSwitch()
	demoAdvancedTypeSwitch()
	demoPatterns()
}
