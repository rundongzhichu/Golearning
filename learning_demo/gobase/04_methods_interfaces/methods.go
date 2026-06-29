// Package methods_interfaces 演示 Go 语言中的方法。
//
// Go 的方法特点：
//   - 方法是绑定了接收者（receiver）的函数
//   - 接收者可以是值类型或指针类型
//   - 只能为同一包内的类型定义方法（不能为内置类型或外部类型添加方法）
//   - 可以通过 type 别名来为「外来」类型添加方法
package main

import (
	"fmt"
	"math"
)

// ==================== 方法定义 ====================

// Vertex 二维点（值类型接收者示例）
type Vertex struct {
	X, Y float64
}

// Abs 值接收者方法：计算到原点的距离
func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Scale 指针接收者方法：缩放坐标
func (v *Vertex) Scale(f float64) {
	v.X *= f
	v.Y *= f
}

// Reset 指针接收者方法：重置为零
func (v *Vertex) Reset() {
	v.X = 0
	v.Y = 0
}

func demoMethodBasics() {
	fmt.Println("=== 方法基础 ===")

	v := Vertex{3, 4}
	fmt.Printf("Vertex%v.Abs() = %.2f\n", v, v.Abs())

	// 指针接收者：修改原值
	v.Scale(2)
	fmt.Printf("Scale(2) 后: %v, Abs() = %.2f\n", v, v.Abs())

	// 值类型自动取地址调用指针接收者方法
	v2 := Vertex{1, 2}
	v2.Reset() // Go 自动转为 (&v2).Reset()
	fmt.Printf("Reset 后: %v\n", v2)
}

// ==================== 接收者选择 ====================

// Counter 计数器（指针接收者的典型用法）
type Counter struct {
	count int
}

// Increment 必须使用指针接收者：需要修改接收者的状态
func (c *Counter) Increment() {
	c.count++
}

// IncrementBy 指针接收者
func (c *Counter) IncrementBy(n int) {
	c.count += n
}

// Value 可以使用值接收者：只读操作
func (c Counter) Value() int {
	return c.count
}

// String 实现 fmt.Stringer 接口
func (c Counter) String() string {
	return fmt.Sprintf("Counter(%d)", c.count)
}

func demoReceiverChoice() {
	fmt.Println("\n=== 接收者选择规则 ===")

	c := &Counter{}
	c.Increment()
	c.IncrementBy(5)
	fmt.Printf("%s.Value() = %d\n", c, c.Value())

	// 一般规则：
	// 1. 需要修改接收者 → 指针接收者
	// 2. 接收者是大结构体 → 指针接收者（避免拷贝）
	// 3. 一致性：如果该类型有任何一个方法使用指针接收者，其余方法也应使用指针接收者
	// 4. 接收者是 map/func/chan → 值接收者（它们本身是引用类型）
}

// ==================== 无法定义方法的情况 ====================

// 不能为非本地类型定义方法
// func (s string) Reverse() string { ... }  // 编译错误！

// 解决方法：使用类型别名
type MyString string

func (s MyString) Reverse() string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func demoTypeAliasMethod() {
	fmt.Println("\n=== 类型别名方法 ===")

	s := MyString("Hello, 世界")
	fmt.Printf("Original: %s\n", s)
	fmt.Printf("Reversed: %s\n", s.Reverse())
}

// ==================== 方法与函数的区别 ====================

func demoMethodVsFunction() {
	fmt.Println("\n=== 方法 vs 函数 ===")

	// 函数
	add := func(a, b int) int { return a + b }
	fmt.Printf("函数: add(1,2)=%d\n", add(1, 2))

	// 方法可以当作函数使用
	v := Vertex{3, 4}
	// v.Abs 是一个方法值（method value）
	distance := v.Abs
	fmt.Printf("方法值: distance()=%.2f\n", distance())

	// 方法表达式（method expression）
	// func (v Vertex) Abs() → func(Vertex) float64
	absFunc := Vertex.Abs
	fmt.Printf("方法表达式: absFunc(v)=%.2f\n", absFunc(v))

	// 指针接收者的方法表达式
	// func (v *Vertex) Scale(f float64) → func(*Vertex, float64)
	scaleFunc := (*Vertex).Scale
	v2 := Vertex{2, 3}
	scaleFunc(&v2, 3)
	fmt.Printf("指针方法表达式: scale后 v2=%v\n", v2)
}

// ==================== 嵌入类型的方法提升 ====================

type Person struct {
	Name string
}

func (p Person) Greet() string {
	return "Hello, I'm " + p.Name
}

type Employee struct {
	Person  // 嵌入 Person
	Company string
}

func demoMethodPromotion() {
	fmt.Println("\n=== 嵌入类型的方法提升 ===")

	e := Employee{
		Person:  Person{Name: "Alice"},
		Company: "ACME Corp",
	}

	// 直接访问提升后的方法
	fmt.Println(e.Greet())

	// 也可以显式调用
	fmt.Println(e.Person.Greet())
}

func main() {
	demoMethodBasics()
	demoReceiverChoice()
	demoTypeAliasMethod()
	demoMethodVsFunction()
	demoMethodPromotion()
}
