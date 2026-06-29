// Package generics 演示 Go 语言中的泛型（Go 1.18+）。
//
// 泛型的特点：
//   - 使用类型参数（type parameters）实现泛型
//   - 类型约束（type constraints）：限制类型参数的范围
//   - 接口可以约束类型参数
//   - 类型推导：很多情况下可以省略显式的类型参数
//   - ~T 近似约束：允许底层类型为 T 的所有类型
//   - comparable 约束：允许 == 和 != 操作
//   - 支持泛型函数、泛型类型（struct）、泛型方法
//   - Go 1.21+：slices、maps、cmp 标准库包
package main

import (
	"cmp"
	"fmt"
	"slices"
)

// ==================== 泛型函数 ====================

// Min 泛型最小值函数
// [T cmp.Ordered] 是类型参数列表，cmp.Ordered 是约束
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max 泛型最大值函数
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// PrintSlice 泛型打印切片函数
func PrintSlice[T any](s []T) {
	for i, v := range s {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}

func demoGenericFunctions() {
	fmt.Println("=== 泛型函数 ===")

	// 类型自动推导
	fmt.Printf("Min(3, 5) = %d\n", Min(3, 5))
	fmt.Printf("Max(3.14, 2.71) = %f\n", Max(3.14, 2.71))
	fmt.Printf("Min('hello', 'world') = %s\n", Min("hello", "world"))

	fmt.Print("PrintSlice: ")
	PrintSlice([]int{1, 2, 3, 4})
	fmt.Print("PrintSlice: ")
	PrintSlice([]string{"Go", "Rust", "Python"})
}

// ==================== 泛型类型 ====================

// Stack 泛型栈
type Stack[T any] struct {
	items []T
}

// Push 入栈
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop 出栈
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T // 零值
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

// Len 栈长度
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Peek 查看栈顶
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func demoGenericType() {
	fmt.Println("\n=== 泛型类型 ===")

	// 整数栈
	intStack := &Stack[int]{}
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)
	fmt.Printf("intStack: len=%d\n", intStack.Len())
	if v, ok := intStack.Pop(); ok {
		fmt.Printf("Pop: %d\n", v)
	}

	// 字符串栈
	strStack := &Stack[string]{}
	strStack.Push("hello")
	strStack.Push("world")
	fmt.Printf("strStack: len=%d\n", strStack.Len())
	if v, ok := strStack.Peek(); ok {
		fmt.Printf("Peek: %s\n", v)
	}
}

// ==================== 类型约束 ====================

// Number 自定义数值约束（~ 允许底层类型为 int/float64 的类型）
type Number interface {
	~int | ~int64 | ~float64
}

// Add 泛型加法（使用自定义约束）
func Add[T Number](a, b T) T {
	return a + b
}

// MyInt 自定义整型（底层类型是 int）
type MyInt int

// Stringer 约束：同时满足 fmt.Stringer 和 Number
type StringerNumber interface {
	fmt.Stringer
	Number
}

func demoConstraints() {
	fmt.Println("\n=== 类型约束 ===")

	fmt.Printf("Add(1, 2) = %d\n", Add(1, 2))
	fmt.Printf("Add(3.5, 2.5) = %.1f\n", Add(3.5, 2.5))

	// MyInt 底层是 int，满足 ~int 约束
	var a, b MyInt = 10, 20
	fmt.Printf("Add(MyInt(10), MyInt(20)) = %d\n", Add(a, b))

	// cmp.Ordered 约束包含所有可比较有序类型
	fmt.Printf("cmp.Ordered Min: %v\n", Min('a', 'z'))
}

// ==================== 泛型切片/映射操作 ====================

// Map 泛型 map 操作：对切片每个元素应用函数
func Map[T any, U any](s []T, fn func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Filter 泛型过滤
func Filter[T any](s []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce 泛型聚合
func Reduce[T any, U any](s []T, initial U, fn func(U, T) U) U {
	result := initial
	for _, v := range s {
		result = fn(result, v)
	}
	return result
}

func demoGenericCollections() {
	fmt.Println("\n=== 泛型集合操作 ===")

	nums := []int{1, 2, 3, 4, 5}

	// Map: 平方
	squares := Map(nums, func(n int) int { return n * n })
	fmt.Printf("Map (平方): %v\n", squares)

	// Filter: 偶数
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Filter (偶数): %v\n", evens)

	// Reduce: 求和
	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Printf("Reduce (求和): %d\n", sum)

	// 链式调用
	result := Reduce(
		Filter(
			Map(nums, func(n int) int { return n * n }),
			func(n int) bool { return n > 10 },
		),
		0,
		func(acc, n int) int { return acc + n },
	)
	fmt.Printf("Map→Filter→Reduce: %d\n", result)
}

// ==================== comparable 约束 ====================

// MapKeys 获取 map 的 keys（comparable 约束）
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func demoComparableConstraint() {
	fmt.Println("\n=== comparable 约束 ===")

	m := map[string]int{"Go": 2009, "Python": 1991, "Rust": 2010}
	keys := MapKeys(m)
	fmt.Printf("keys: %v\n", keys)

	// comparable 约束允许 == 和 != 操作
	// 可作为 map 的 key 的类型都满足 comparable
	fmt.Printf("slices.Contains: %v\n", slices.Contains(keys, "Go"))
}

func main() {
	demoGenericFunctions()
	demoGenericType()
	demoConstraints()
	demoGenericCollections()
	demoComparableConstraint()
}
