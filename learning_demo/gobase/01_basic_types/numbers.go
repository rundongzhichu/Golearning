// Package basic_types 演示 Go 语言中的数字类型（整数与浮点数）。
//
// Go 的数字类型分为：
//   - 整数：int8, int16, int32, int64, int（平台相关，64位平台等同于 int64）
//   - 无符号整数：uint8(byte), uint16, uint32, uint64, uint(rune)
//   - 浮点数：float32, float64
//   - 复数：complex64, complex128
//   - uintptr：用于存储指针的整数类型
//
// 零值：所有数字类型的零值都是 0。
package main

import (
	"fmt"
	"math"
	"strconv"
)

// ==================== 整数 ====================

// demoIntegerTypes 演示各种整数类型的声明、默认值与运算。
func demoIntegerTypes() {
	// 显式类型声明
	var a int8 = 127                  // int8 范围：-128 ~ 127
	var b int16 = 32767               // int16 范围：-32768 ~ 32767
	var c int32 = 2147483647          // int32（即 rune）范围：约 ±21 亿
	var d int64 = 9223372036854775807 // int64 范围：约 ±9.2e18

	// int 和 uint 的大小依赖于平台（64 位平台上 64 位）
	var e int = 100
	var f uint = 200

	// 类型推断（默认推断为 int）
	g := 42

	// byte 是 uint8 的别名
	var h byte = 255

	fmt.Println("=== 整数类型 ===")
	fmt.Printf("int8:  %d\n", a)
	fmt.Printf("int16: %d\n", b)
	fmt.Printf("int32: %d\n", c)
	fmt.Printf("int64: %d\n", d)
	fmt.Printf("int:   %d\n", e)
	fmt.Printf("uint:  %d\n", f)
	fmt.Printf("推断:  %d (类型: %T)\n", g, g)
	fmt.Printf("byte:  %d (类型: %T)\n", h, h)

	// 整数运算
	sum := a + 10
	diff := d - 1
	prod := e * 2
	quot := f / 3      // 整数除法：向下取整
	remainder := e % 3 // 取余
	fmt.Printf("\n运算: sum=%d, diff=%d, prod=%d, quot=%d, remainder=%d\n",
		sum, diff, prod, quot, remainder)

	// 类型转换（Go 不允许隐式类型转换）
	var i int64 = 100
	j := int(i)   // int64 → int
	k := int32(i) // int64 → int32
	fmt.Printf("类型转换: int64(%d) → int(%d) → int32(%d)\n", i, j, k)
}

// ==================== 浮点数 ====================

// demoFloatTypes 演示浮点数类型的声明、精度及 math 包的使用。
func demoFloatTypes() {
	// float32 和 float64
	var f32 float32 = 3.141592653589793
	var f64 float64 = 3.141592653589793

	fmt.Println("\n=== 浮点数类型 ===")
	fmt.Printf("float32: %.15f (精度损失)\n", f32)
	fmt.Printf("float64: %.15f (高精度)\n", f64)

	// 科学记数法
	big := 1.5e10   // 1.5 * 10^10
	small := 2.3e-4 // 2.3 * 10^-4
	fmt.Printf("科学记数: big=%.0f, small=%.6f\n", big, small)

	// math 包：数学函数
	fmt.Println("\n--- math 包 ---")
	fmt.Printf("math.Pi:       %v\n", math.Pi)
	fmt.Printf("math.MaxFloat64: %v\n", math.MaxFloat64)
	fmt.Printf("math.Sqrt(16): %v\n", math.Sqrt(16))
	fmt.Printf("math.Pow(2,10): %v\n", math.Pow(2, 10))
	fmt.Printf("math.Abs(-5):  %v\n", math.Abs(-5))
	fmt.Printf("math.Ceil(3.14):  %v\n", math.Ceil(3.14))
	fmt.Printf("math.Floor(3.14): %v\n", math.Floor(3.14))
	fmt.Printf("math.Round(3.54): %v\n", math.Round(3.54))

	// 浮点数比较：使用差值比较而非 ==
	a, b := 0.1+0.2, 0.3
	epsilon := 1e-9
	fmt.Printf("\n0.1+0.2 == 0.3? %v\n", a == b)
	fmt.Printf("差值 < epsilon? %v (差值: %v)\n",
		math.Abs(a-b) < epsilon, math.Abs(a-b))
}

// ==================== 复数 ====================

// demoComplex 演示复数的声明与运算。
func demoComplex() {
	// 复数声明
	c1 := complex(3, 4) // 3 + 4i
	c2 := 1 + 2i        // 直接字面量

	fmt.Println("\n=== 复数类型 ===")
	fmt.Printf("c1:       %v\n", c1)
	fmt.Printf("c2:       %v\n", c2)
	fmt.Printf("实部:     %v\n", real(c1))
	fmt.Printf("虚部:     %v\n", imag(c1))
	fmt.Printf("加法:     %v\n", c1+c2)
	fmt.Printf("乘法:     %v\n", c1*c2)
	fmt.Printf("模(abs):  %v\n", cmplxAbs(c1))
}

// cmplxAbs 计算复数的模（math/cmplx.Abs 的演示替代）。
func cmplxAbs(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}

// ==================== 字符串与数字互转 ====================

// demoStringNumberConversion 演示字符串与数字之间的转换（strconv 包）。
func demoStringNumberConversion() {
	fmt.Println("\n=== 字符串与数字转换 ===")

	// 数字 → 字符串
	s1 := strconv.Itoa(42)                         // int → string
	s2 := strconv.FormatInt(255, 16)               // int64 → 16进制字符串
	s3 := strconv.FormatFloat(3.14159, 'f', 4, 64) // float64 → 保留4位小数
	fmt.Printf("Itoa:        %q\n", s1)
	fmt.Printf("FormatInt:   %q\n", s2)
	fmt.Printf("FormatFloat: %q\n", s3)

	// 字符串 → 数字
	n1, _ := strconv.Atoi("42")                // string → int
	n2, _ := strconv.ParseInt("FF", 16, 64)    // 16进制字符串 → int64
	n3, _ := strconv.ParseFloat("3.14159", 64) // string → float64
	fmt.Printf("Atoi:        %d\n", n1)
	fmt.Printf("ParseInt:    %d\n", n2)
	fmt.Printf("ParseFloat:  %f\n", n3)

	// 使用 fmt.Sprintf
	s4 := fmt.Sprintf("%d", 100)
	fmt.Printf("Sprintf:     %q\n", s4)
}

func main() {
	demoIntegerTypes()
	demoFloatTypes()
	demoComplex()
	demoStringNumberConversion()
}
