// Package basic_types 演示 Go 语言中的字符串类型及相关操作。
//
// Go 的字符串特点：
//   - string 是 Go 的内置类型，零值为空字符串 ""
//   - 字符串是不可变的（immutable）：一旦创建，内容不可修改
//   - 字符串底层是 []byte，但使用 UTF-8 编码
//   - 可以用 “（反引号）创建原始字符串（raw string），忽略转义
//   - 可以用 "" 创建解释型字符串，支持 \n、\t 等转义
//   - strings 包提供丰富的字符串操作函数
//   - strconv 包提供字符串与基本类型间的转换
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// demoStringDeclaration 演示字符串的声明方式。
func demoStringDeclaration() {
	// 显式声明
	var s1 string = "Hello, Go!"
	// 短变量声明
	s2 := "你好，世界"
	// 零值
	var s3 string
	// 原始字符串（raw string）：不转义，适合正则、路径、多行文本
	s4 := `第一行
第二行
第三行\t\n 这里的转义符无效`

	fmt.Println("=== 字符串声明 ===")
	fmt.Printf("显式声明:   %s\n", s1)
	fmt.Printf("短变量声明: %s\n", s2)
	fmt.Printf("零值:       %q (空字符串)\n", s3)
	fmt.Printf("原始字符串:\n%s\n", s4)

	// 字符串是不可变的，不能直接修改某个字符
	// s1[0] = 'h'  // 编译错误：cannot assign to s1[0]
	fmt.Println("注意: Go 的字符串是不可变的 (immutable)")
}

// demoStringLength 演示字符串长度计算。
func demoStringLength() {
	s := "Hello, 世界"

	fmt.Println("\n=== 字符串长度 ===")
	// len() 返回字节数，不是字符数
	fmt.Printf("字符串:     %s\n", s)
	fmt.Printf("字节数(len): %d\n", len(s))
	// utf8.RuneCountInString() 返回 Unicode 字符数（rune 数）
	fmt.Printf("字符数:      %d\n", utf8.RuneCountInString(s))
}

// demoStringIteration 演示字符串的遍历方式。
func demoStringIteration() {
	s := "Hello, 世界"

	fmt.Println("\n=== 字符串遍历 ===")

	// 方式1：for range 遍历 rune（字符），推荐！
	fmt.Print("for range (按字符/rune): ")
	for i, r := range s {
		fmt.Printf("[%d:%c] ", i, r)
	}
	fmt.Println()

	// 方式2：for i 按字节遍历（不推荐用于多字节字符）
	fmt.Print("for i (按字节):          ")
	for i := 0; i < len(s); i++ {
		fmt.Printf("[%d:%c] ", i, s[i])
	}
	fmt.Println()
}

// demoStringOperations 演示 strings 包的常用函数。
func demoStringOperations() {
	fmt.Println("\n=== 字符串操作 (strings 包) ===")

	s := "Hello, Go!"

	// 包含
	fmt.Println("Contains 'Go':", strings.Contains(s, "Go"))
	fmt.Println("Contains 'Java':", strings.Contains(s, "Java"))

	// 前缀/后缀
	fmt.Println("HasPrefix 'He':", strings.HasPrefix(s, "He"))
	fmt.Println("HasSuffix '!':", strings.HasSuffix(s, "!"))

	// 查找
	fmt.Println("Index 'Go':", strings.Index(s, "Go"))
	fmt.Println("LastIndex 'o':", strings.LastIndex(s, "o"))

	// 分割与合并
	parts := strings.Split("a,b,c", ",")
	fmt.Println("Split:", parts)
	fmt.Println("Join:", strings.Join(parts, "-"))

	// 大小写转换
	fmt.Println("ToUpper:", strings.ToUpper(s))
	fmt.Println("ToLower:", strings.ToLower(s))

	// 替换
	fmt.Println("Replace:", strings.Replace(s, "o", "0", 1))    // 替换1次
	fmt.Println("ReplaceAll:", strings.ReplaceAll(s, "o", "0")) // 替换全部

	// 去除空白（TrimSpace）、指定字符（Trim）
	fmt.Println("TrimSpace:", strings.TrimSpace("  hello  "))
	fmt.Println("Trim:", strings.Trim("!!!hello!!!", "!"))

	// 重复
	fmt.Println("Repeat:", strings.Repeat("Go! ", 3))

	// 比较（区分大小写）
	fmt.Println("EqualFold:", strings.EqualFold("Go", "go")) // 不区分大小写比较

	// Builder：高效拼接字符串
	var builder strings.Builder
	builder.WriteString("Hello")
	builder.WriteString(", ")
	builder.WriteString("World!")
	fmt.Println("Builder:", builder.String())
}

// demoStringConversion 演示字符串与其他类型的互转。
func demoStringConversion() {
	fmt.Println("\n=== 字符串转换 ===")

	// []byte ↔ string（会产生内存拷贝，注意性能）
	s := "hello"
	b := []byte(s)  // string → []byte
	s2 := string(b) // []byte → string
	fmt.Printf("string→[]byte: %v\n", b)
	fmt.Printf("[]byte→string: %s\n", s2)

	// []rune ↔ string（适合操作 Unicode 字符）
	runes := []rune(s)
	runes[0] = 'H' // 修改第一个字符
	s3 := string(runes)
	fmt.Printf("[]rune→string (修改后): %s\n", s3)

	// 单字符与 rune
	var r rune = '中'
	fmt.Printf("rune: %c (Unicode: U+%04X)\n", r, r)
}

func main() {
	demoStringDeclaration()
	demoStringLength()
	demoStringIteration()
	demoStringOperations()
	demoStringConversion()
}
