// Package methods_interfaces 演示 Go 语言中的空接口（interface{} / any）。
//
// 空接口的特点：
//   - 可以存储任意类型的值（类似其他语言的 Object/void*）
//   - Go 1.18+ 引入了 any 作为 interface{} 的别名
//   - 常用于 JSON 解析、fmt.Println、容器类型等场景
//   - 从空接口取出值需要类型断言
package main

import (
	"encoding/json"
	"fmt"
)

// ==================== 空接口基础 ====================

func demoEmptyInterfaceBasics() {
	fmt.Println("=== 空接口基础 ===")

	// Go 1.18+ 可以用 any 替代 interface{}
	var v any

	// 可以存储任意类型
	v = 42
	fmt.Printf("int: %v (type: %T)\n", v, v)

	v = "hello"
	fmt.Printf("string: %v (type: %T)\n", v, v)

	v = struct{ Name string }{"Alice"}
	fmt.Printf("struct: %v (type: %T)\n", v, v)

	v = []any{1, "two", 3.0}
	fmt.Printf("slice: %v (type: %T)\n", v, v)
}

// ==================== 空接口作为容器 ====================

// PrintAll 接收任意数量、任意类型的参数（类似 fmt.Println）
func PrintAll(args ...any) {
	for i, arg := range args {
		fmt.Printf("[%d] %v (type: %T)\n", i, arg, arg)
	}
}

// Stack 使用空接口的通用栈
type Stack struct {
	items []any
}

// Push 入栈
func (s *Stack) Push(item any) {
	s.items = append(s.items, item)
}

// Pop 出栈（需要类型断言）
func (s *Stack) Pop() any {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func demoContainer() {
	fmt.Println("\n=== 空接口容器 ===")

	// PrintAll 接收任意类型
	PrintAll(42, "hello", 3.14, true)

	// Stack 示例
	stack := &Stack{}
	stack.Push(1)
	stack.Push("middle")
	stack.Push(3.14)

	fmt.Println("\nStack 出栈:")
	for i := 0; i < 3; i++ {
		item := stack.Pop()
		fmt.Printf("  %v (type: %T)\n", item, item)
	}
}

// ==================== JSON 解析（空接口典型应用） ====================

func demoJSONWithEmptyInterface() {
	fmt.Println("\n=== JSON 与空接口 ===")

	// 解析 JSON 到空接口（结构未知时使用）
	jsonStr := `{
		"name": "Alice",
		"age": 30,
		"hobbies": ["reading", "coding"],
		"address": {
			"city": "Beijing",
			"zip": "100000"
		}
	}`

	var result any
	json.Unmarshal([]byte(jsonStr), &result)

	// 通过类型断言访问嵌套数据
	m := result.(map[string]any)
	fmt.Printf("name: %v\n", m["name"])
	fmt.Printf("age: %v\n", m["age"])

	hobbies := m["hobbies"].([]any)
	fmt.Printf("hobbies: %v\n", hobbies)

	address := m["address"].(map[string]any)
	fmt.Printf("city: %v\n", address["city"])
}

// ==================== 空接口与 nil ====================

type MyError struct{}

func (e *MyError) Error() string { return "my error" }

// getError 返回一个非 nil 接口值（接口有类型信息，值为 nil）
func getError() error {
	var e *MyError // nil 指针
	return e       // 返回的 error 接口不是 nil！（接口持有 *MyError 类型信息）
}

func demoNilInterface() {
	fmt.Println("\n=== 空接口与 nil 陷阱 ===")

	err := getError()
	fmt.Printf("err == nil? %v\n", err == nil) // false！
	fmt.Printf("err 的值: %v, 类型: %T\n", err, err)

	// 正确的做法：返回显式的 nil
	correctGetError := func() error {
		return nil // 直接返回 nil
	}
	err2 := correctGetError()
	fmt.Printf("correct: err == nil? %v\n", err2 == nil) // true
}

// ==================== Go 1.18+ any 类型 ====================

func demoAnyType() {
	fmt.Println("\n=== Go 1.18+ any 类型 ===")

	// any 是 interface{} 的别名，完全等价
	var x any = "hello"
	var y interface{} = "hello"

	fmt.Printf("x == y? %v\n", x == y)
	fmt.Printf("x 类型: %T, y 类型: %T\n", x, y)

	// any 在泛型中也可以使用
	printAny := func(v any) {
		fmt.Printf("any value: %v\n", v)
	}
	printAny(42)
	printAny("world")
}

func main() {
	demoEmptyInterfaceBasics()
	demoContainer()
	demoJSONWithEmptyInterface()
	demoNilInterface()
	demoAnyType()
}
