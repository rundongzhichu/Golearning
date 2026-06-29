// Package composite_types 演示 Go 语言中的切片类型（slice）。
//
// 切片是 Go 中最常用的数据结构之一，特点如下：
//   - 动态数组：长度可变，基于数组实现
//   - 引用类型：赋值和传参只复制 header（指针+长度+容量），不复制底层数组
//   - 零值为 nil
//   - 底层结构：ptr（指向底层数组）、len（长度）、cap（容量）
//   - 使用 make() 或字面量创建
//   - append() 追加元素，自动扩容
//   - copy() 复制元素
package main

import (
	"fmt"
	"slices"
)

// ==================== 切片创建 ====================

func demoSliceCreation() {
	fmt.Println("=== 切片创建 ===")

	// 方式1：声明 nil 切片（零值）
	var s1 []int
	fmt.Printf("nil切片: %v, len=%d, cap=%d, isNil=%v\n", s1, len(s1), cap(s1), s1 == nil)

	// 方式2：使用字面量
	s2 := []int{1, 2, 3, 4, 5}
	fmt.Printf("字面量: %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	// 方式3：使用 make([]T, len, cap)
	s3 := make([]int, 3, 5) // len=3, cap=5
	fmt.Printf("make:   %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	// 方式4：从数组或切片创建（切片表达式）
	arr := [5]int{10, 20, 30, 40, 50}
	s4 := arr[1:4] // 从 arr[1] 到 arr[3]，左闭右开
	fmt.Printf("切片表达式 arr[1:4]: %v, len=%d, cap=%d\n", s4, len(s4), cap(s4))

	// 方式5：空切片（非 nil）
	s5 := make([]int, 0)
	fmt.Printf("空切片: %v, len=%d, cap=%d, isNil=%v\n", s5, len(s5), cap(s5), s5 == nil)
}

// ==================== 切片操作 ====================

func demoSliceOperations() {
	fmt.Println("\n=== 切片操作 ===")

	s := []int{1, 2, 3}

	// append：追加元素
	s = append(s, 4)
	s = append(s, 5, 6, 7)        // 追加多个
	s = append(s, []int{8, 9}...) // 展开另一个切片（使用 ...）
	fmt.Printf("append 后: %v, len=%d, cap=%d\n", s, len(s), cap(s))

	// copy：复制元素（返回复制的元素个数）
	src := []int{100, 200, 300}
	dst := make([]int, 2)
	n := copy(dst, src) // 复制 min(len(dst), len(src)) 个
	fmt.Printf("copy: src=%v, dst=%v, copied=%d\n", src, dst, n)

	// 切片表达式 [low:high] 和 [low:high:max]
	base := []int{1, 2, 3, 4, 5}
	a := base[1:3]   // [2,3], len=2, cap=4 (从base[1]到底层数组末尾)
	b := base[1:3:3] // [2,3], len=2, cap=2 (限制容量到3-1=2)
	fmt.Printf("base[1:3]:   %v, len=%d, cap=%d\n", a, len(a), cap(a))
	fmt.Printf("base[1:3:3]: %v, len=%d, cap=%d\n", b, len(b), cap(b))

	// 删除元素（使用 append + 切片操作）
	s = []int{1, 2, 3, 4, 5}
	// 删除索引 2 的元素（值为 3）
	s = append(s[:2], s[3:]...)
	fmt.Printf("删除索引2: %v\n", s)

	// 插入元素（使用 append + 切片操作）
	s = []int{1, 2, 4, 5}
	// 在索引 2 处插入 3
	s = append(s[:2], append([]int{3}, s[2:]...)...)
	fmt.Printf("插入索引2: %v\n", s)

	// 清空切片（保留底层数组容量）
	s = s[:0]
	fmt.Printf("清空后: %v, len=%d, cap=%d\n", s, len(s), cap(s))
}

// ==================== 切片扩容 ====================

func demoSliceGrow() {
	fmt.Println("\n=== 切片扩容机制 ===")

	s := make([]int, 0)
	fmt.Println("append 过程中的 len/cap 变化：")
	for i := 1; i <= 20; i++ {
		s = append(s, i)
		fmt.Printf("  append(%2d): len=%2d, cap=%2d\n", i, len(s), cap(s))
	}
	// 扩容策略（Go 1.18+）：
	// - 小切片（cap < 256）：每次扩容约 2 倍
	// - 大切片（cap >= 256）：使用公式 newcap += (newcap + 3*256) / 4
}

// ==================== 切片底层共享 ====================

func demoSliceShare() {
	fmt.Println("\n=== 切片底层数组共享 ===")

	base := []int{1, 2, 3, 4, 5}
	sub1 := base[1:4] // [2,3,4], 共享底层数组
	sub2 := base[2:5] // [3,4,5], 共享底层数组

	fmt.Printf("初始: base=%v, sub1=%v, sub2=%v\n", base, sub1, sub2)

	// 修改 sub1 会影响 base 和 sub2（共享底层数组）
	sub1[1] = 999
	fmt.Printf("修改 sub1[1] 后: base=%v, sub1=%v, sub2=%v\n", base, sub1, sub2)

	// 当 append 超过 cap 时，会分配新的底层数组（脱离共享）
	sub1 = append(sub1, 6, 7, 8) // 超过了 cap
	sub1[0] = 888                // 现在不会影响 base
	fmt.Printf("append 超出容量后 sub1[0]=888: base=%v, sub1=%v\n", base, sub1)
}

// ==================== 切片遍历与工具函数 ====================

func demoSliceIteration() {
	fmt.Println("\n=== 切片遍历与工具函数 ===")

	s := []string{"Go", "Python", "Rust", "JavaScript"}

	// for range 遍历（Go 1.22+ 支持 range over int）
	fmt.Print("for range: ")
	for i, v := range s {
		fmt.Printf("[%d]=%s ", i, v)
	}
	fmt.Println()

	// 传统 for 循环
	fmt.Print("经典 for:  ")
	for i := 0; i < len(s); i++ {
		fmt.Printf("%s ", s[i])
	}
	fmt.Println()

	// slices 包（Go 1.21+ 标准库）
	fmt.Println("\n--- slices 包 ---")
	fmt.Printf("Contains 'Go': %v\n", slices.Contains(s, "Go"))
	fmt.Printf("Index 'Rust': %d\n", slices.Index(s, "Rust"))

	// 排序
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	slices.Sort(nums)
	fmt.Printf("排序后: %v\n", nums)

	// 最大值/最小值
	fmt.Printf("Max: %d, Min: %d\n", slices.Max(nums), slices.Min(nums))

	// 比较两个切片
	fmt.Printf("Equal: %v\n", slices.Equal(nums, []int{1, 1, 2, 3, 4, 5, 6, 9}))
}

func main() {
	demoSliceCreation()
	demoSliceOperations()
	demoSliceGrow()
	demoSliceShare()
	demoSliceIteration()
}
