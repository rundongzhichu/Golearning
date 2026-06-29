// Package composite_types 演示 Go 语言中的 map 类型（映射/字典）。
//
// map 的特点：
//   - 键值对集合：无序的 key-value 结构
//   - 引用类型：赋值和传参只复制 header
//   - 零值为 nil；nil map 不能写入，但可以读取
//   - 使用 make() 或字面量创建
//   - key 必须是可比较类型（不能是 slice/map/function）
//   - 访问不存在的 key 会返回零值，使用 comma-ok 语法判断是否存在
//   - delete() 函数用于删除键值对
package main

import (
	"fmt"
	"maps"
)

// ==================== map 创建 ====================

func demoMapCreation() {
	fmt.Println("=== map 创建 ===")

	// 方式1：声明 nil map（不能写入！）
	var m1 map[string]int
	fmt.Printf("nil map: %v, len=%d, isNil=%v\n", m1, len(m1), m1 == nil)
	// m1["key"] = 1  // 运行时 panic: assignment to entry in nil map

	// 方式2：使用 make
	m2 := make(map[string]int)
	fmt.Printf("make: %v, len=%d\n", m2, len(m2))

	// 方式3：make 指定初始容量（减少扩容）
	m3 := make(map[string]int, 100)
	fmt.Printf("make with cap: %v, len=%d\n", m3, len(m3))

	// 方式4：字面量初始化
	m4 := map[string]int{
		"Go":     2009,
		"Python": 1991,
		"Rust":   2010,
	}
	fmt.Printf("字面量: %v\n", m4)

	// 方式5：空 map 字面量
	m5 := map[string]int{}
	fmt.Printf("空字面量: %v, len=%d, isNil=%v\n", m5, len(m5), m5 == nil)
}

// ==================== map 基本操作 ====================

func demoMapOperations() {
	fmt.Println("\n=== map 基本操作 ===")

	m := make(map[string]int)

	// 写入
	m["apple"] = 5
	m["banana"] = 3
	m["orange"] = 8

	// 读取
	fmt.Printf("apple=%d, banana=%d\n", m["apple"], m["banana"])

	// 读取不存在的 key：返回零值（不会报错）
	fmt.Printf("不存在的 key 'grape'=%d (零值)\n", m["grape"])

	// 判断 key 是否存在：comma-ok 语法
	value, ok := m["apple"]
	fmt.Printf("'apple' exists: %v, value=%d\n", ok, value)

	value, ok = m["grape"]
	fmt.Printf("'grape' exists: %v, value=%d\n", ok, value)

	// 删除
	delete(m, "banana")
	fmt.Printf("删除 'banana' 后: %v\n", m)

	// 删除不存在的 key 也不会报错
	delete(m, "not_exist")

	// 长度
	fmt.Printf("len(m)=%d\n", len(m))

	// 清空 map（重新 make 或逐条删除）
	// 方式1：重新 make（推荐，GC 会回收旧的）
	m = make(map[string]int)
	fmt.Printf("重新 make 后: %v, len=%d\n", m, len(m))
}

// ==================== map 遍历 ====================

func demoMapIteration() {
	fmt.Println("\n=== map 遍历 ===")

	m := map[string]int{
		"Go":     2009,
		"Python": 1991,
		"Rust":   2010,
		"C":      1972,
	}

	// 遍历 key-value（注意：顺序不确定！）
	fmt.Println("for range (顺序不确定):")
	for k, v := range m {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// 只遍历 key
	fmt.Print("只遍历 key: ")
	for k := range m {
		fmt.Printf("%s ", k)
	}
	fmt.Println()

	// 只遍历 value
	fmt.Print("只遍历 value: ")
	for _, v := range m {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}

// ==================== map 进阶 ====================

func demoMapAdvanced() {
	fmt.Println("\n=== map 进阶用法 ===")

	// map 作为集合（set）：使用 map[key]bool 或 map[key]struct{}
	// struct{} 不占内存，更推荐
	set := make(map[string]struct{})
	set["Go"] = struct{}{}
	set["Python"] = struct{}{}
	set["Go"] = struct{}{} // 重复插入无影响
	fmt.Printf("set: len=%d\n", len(set))
	if _, ok := set["Go"]; ok {
		fmt.Println("Go 在集合中")
	}

	// 嵌套 map
	nested := make(map[string]map[string]int)
	nested["fruits"] = map[string]int{"apple": 5}
	nested["vegetables"] = map[string]int{"carrot": 3}
	fmt.Printf("嵌套 map: %v\n", nested)
	fmt.Printf("fruits.apple = %d\n", nested["fruits"]["apple"])

	// map 的 value 可以是函数
	funcMap := map[string]func(int, int) int{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}
	fmt.Printf("add(3,4)=%d, sub(10,3)=%d\n",
		funcMap["add"](3, 4), funcMap["sub"](10, 3))

	// 统计单词频率
	words := []string{"go", "python", "go", "rust", "go", "python"}
	counter := make(map[string]int)
	for _, w := range words {
		counter[w]++
	}
	fmt.Printf("词频统计: %v\n", counter)

	// maps 包（Go 1.21+）
	fmt.Println("\n--- maps 包 ---")
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := maps.Clone(m1) // 浅拷贝
	fmt.Printf("Clone: %v, Equal: %v\n", m2, maps.Equal(m1, m2))
}

func main() {
	demoMapCreation()
	demoMapOperations()
	demoMapIteration()
	demoMapAdvanced()
}
