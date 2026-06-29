// Package basic_types 演示 Go 语言中的布尔类型。
//
// Go 的布尔类型只有两个值：true 和 false，零值为 false。
// 与其他语言不同，Go 的 bool 和 int 不能互相转换（如 1 不能代表 true）。
package main

import "fmt"

func main() {
	// ==================== 布尔声明与零值 ====================
	var b1 bool // 零值为 false
	var b2 bool = true
	b3 := false // 短变量声明
	b4 := 1 < 2 // 通过比较表达式

	fmt.Println("=== 布尔声明 ===")
	fmt.Printf("零值:  %v\n", b1)
	fmt.Printf("true:  %v\n", b2)
	fmt.Printf("false: %v\n", b3)
	fmt.Printf("1 < 2: %v\n", b4)

	// ==================== 逻辑运算符 ====================
	fmt.Println("\n=== 逻辑运算符 ===")
	fmt.Println("&& (AND): true && true  =", true && true)
	fmt.Println("&& (AND): true && false =", true && false)
	fmt.Println("|| (OR):  true || false =", true || false)
	fmt.Println("|| (OR):  false|| false =", false || false)
	fmt.Println("!  (NOT): !true          =", !true)
	fmt.Println("!  (NOT): !false         =", !false)

	// 短路求值：如果左边已经能决定结果，右边不会被求值
	fmt.Println("\n--- 短路求值 ---")
	a := 0
	// 左边为 false，右边不会被计算（a 保持 0）
	if false && (func() bool { a = 100; return true }()) {
		fmt.Println("不会执行")
	}
	fmt.Printf("短路 AND 后 a=%d (未被修改)\n", a)

	// 左边为 true，右边不会被计算
	if true || (func() bool { a = 100; return false }()) {
		fmt.Printf("短路 OR 后 a=%d (未被修改)\n", a)
	}

	// ==================== 比较运算符 ====================
	fmt.Println("\n=== 比较运算符 ===")
	fmt.Println("==: 5 == 5  =", 5 == 5)
	fmt.Println("!=: 5 != 3  =", 5 != 3)
	fmt.Println("< : 3 < 5   =", 3 < 5)
	fmt.Println("<=: 5 <= 5  =", 5 <= 5)
	fmt.Println("> : 7 > 3   =", 7 > 3)
	fmt.Println(">=: 5 >= 3  =", 5 >= 3)

	// 字符串也可以比较（按字典序）
	fmt.Println("\"abc\" < \"abd\":", "abc" < "abd")

	// ==================== 条件判断中的常见模式 ====================
	fmt.Println("\n=== 常用模式 ===")

	// if 条件直接使用 bool 变量
	isReady := true
	if isReady {
		fmt.Println("直接使用 bool: isReady 为 true")
	}

	// 多个条件组合
	age, hasLicense := 22, true
	if age >= 18 && hasLicense {
		fmt.Println("组合条件: 可以开车")
	}

	// bool 取反
	closed := false
	if !closed {
		fmt.Println("取反: 通道未关闭")
	}
}
