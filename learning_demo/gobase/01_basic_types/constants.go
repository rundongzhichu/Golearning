// Package basic_types 演示 Go 语言中的常量与 iota 枚举。
//
// 常量的特点：
//   - 使用 const 关键字声明，在编译期确定值
//   - 可以是布尔、数字（整数、浮点、复数）或字符串
//   - 不能取地址（&）——常量不分配内存
//   - iota 是 Go 的枚举生成器，在每个 const 块内从 0 开始递增
//   - 无类型常量具有更高的精度和灵活性
package main

import (
	"fmt"
	"math"
)

// ==================== 常量声明 ====================

// 单行声明
const Pi = 3.1415926535

// 多行声明
const (
	StatusOK       = 200
	StatusNotFound = 404
	StatusError    = 500
)

// 类型化常量
const (
	KB float64 = 1 << 10 // 1024
	MB         = 1 << 20 // 1048576
	GB         = 1 << 30
)

// 无类型常量：没有显式指定类型，可以有更高精度
const BigNumber = 1 << 100 // 超大值，作为无类型常量没问题

func demoConstants() {
	fmt.Println("=== 常量 ===")
	fmt.Printf("Pi:          %v\n", Pi)
	fmt.Printf("StatusOK:    %d\n", StatusOK)
	fmt.Printf("KB:          %.0f\n", KB)
	fmt.Printf("MB:          %.0f\n", MB)
	fmt.Printf("GB:          %.0f\n", GB)

	// 无类型常量可以在使用时自动适配类型
	var f32 float32 = Pi // 自动适配为 float32
	var f64 float64 = Pi // 自动适配为 float64
	fmt.Printf("float32 Pi:  %v\n", f32)
	fmt.Printf("float64 Pi:  %v\n", f64)
}

// ==================== iota 枚举 ====================

// 基础 iota：从 0 开始递增
const (
	Monday    = iota // 0
	Tuesday          // 1 (隐式重复上一行表达式)
	Wednesday        // 2
	Thursday         // 3
	Friday           // 4
	Saturday         // 5
	Sunday           // 6
)

// iota 可参与表达式
const (
	_   = iota             // 跳过 0
	KB2 = 1 << (10 * iota) // 1 << 10 = 1024
	MB2 = 1 << (10 * iota) // 1 << 20
	GB2 = 1 << (10 * iota) // 1 << 30
)

// 多个 iota 在同一行：iota 在同一行内不变
const (
	a, b = iota, iota + 1 // a=0, b=1
	c, d = iota, iota + 1 // c=1, d=2
)

// 使用 iota 创建位掩码（常见模式）
const (
	FlagRead  = 1 << iota // 1 << 0 = 1
	FlagWrite             // 1 << 1 = 2
	FlagExec              // 1 << 2 = 4
)

func demoIota() {
	fmt.Println("\n=== iota 枚举 ===")
	fmt.Printf("Monday=%d, Tuesday=%d, ..., Sunday=%d\n",
		Monday, Tuesday, Sunday)
	fmt.Printf("KB2=%d, MB2=%d, GB2=%d\n", KB2, MB2, GB2)
	fmt.Printf("a=%d, b=%d, c=%d, d=%d\n", a, b, c, d)
	fmt.Printf("FlagRead=%b, FlagWrite=%b, FlagExec=%b\n",
		FlagRead, FlagWrite, FlagExec)

	// 位掩码用法
	permission := FlagRead | FlagWrite
	fmt.Printf("permission(读+写): %b\n", permission)
	fmt.Printf("有读权限? %v\n", permission&FlagRead != 0)
	fmt.Printf("有执行权限? %v\n", permission&FlagExec != 0)
}

// ==================== 常量生成器 ====================

// 无类型常量的灵活性演示
func demoUntypedConstants() {
	fmt.Println("\n=== 无类型常量的灵活性 ===")

	const (
		Year    = 365.2425 // 无类型浮点常量
		MaxUint = ^uint(0) // 无类型整数常量
	)

	// 可以用于需要不同精度的计算
	fmt.Printf("Year:     %v (float64: %f)\n", Year, float64(Year))
	fmt.Printf("MaxUint:  %v (uint64: %d)\n", MaxUint, uint64(MaxUint))

	// 无类型常量可以有任意精度
	const Huge = 1 << 200
	const Tiny = Huge >> 199
	fmt.Printf("Tiny: %d\n", Tiny)

	// math 包中的常量
	fmt.Printf("math.Pi:  %v\n", math.Pi)
	fmt.Printf("math.E:   %v\n", math.E)
}

func main() {
	demoConstants()
	demoIota()
	demoUntypedConstants()
}
