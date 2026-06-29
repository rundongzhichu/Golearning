// Package packages_modules 演示 Go 语言中的包（package）和模块（module）。
//
// Go 的包管理特点：
//   - 每个 Go 文件第一行必须是 package 声明
//   - 同一目录下所有 .go 文件属于同一个包
//   - 包名通常与目录名一致（但也可以不同）
//   - 首字母大写的标识符（函数/类型/变量）可以被其他包访问（导出）
//   - 首字母小写的标识符是包内私有的
//   - main 包的 main() 函数是程序的入口
//   - go.mod 定义模块名和依赖（Go modules）
//
// 导入：
//   - import "fmt"：标准导入
//   - import "net/http"：使用路径最后一段作为包名
//   - import f "fmt"：给包起别名
//   - import . "fmt"：将包内容导入当前命名空间（不推荐）
//   - import _ "net/http/pprof"：仅执行 init 函数（side-effect only）
//
// 本文件展示包内可见性、init 函数、导入声明等概念
package main

import (
	"fmt"
	"math/rand"
	"time"
	// _ "net/http/pprof"  // 仅执行 init：注册 pprof 路由
)

// ==================== 导出 vs 未导出 ====================

// ExportedType 可以被其他包访问（首字母大写）
type ExportedType struct {
	ExportedField   string // 导出字段
	unexportedField string // 包内私有字段
}

// ExportedFunc 导出函数
func ExportedFunc() string {
	return "I'm exported!"
}

// unexportedFunc 包内私有函数
func unexportedFunc() string {
	return "internal only"
}

func demoVisibility() {
	fmt.Println("=== 包可见性 ===")

	t := ExportedType{
		ExportedField:   "visible",
		unexportedField: "hidden",
	}
	fmt.Printf("导出字段: %s\n", t.ExportedField)
	fmt.Printf("私有字段(包内可访问): %s\n", t.unexportedField)
	fmt.Printf("导出函数: %s\n", ExportedFunc())
	fmt.Printf("私有函数(包内可调用): %s\n", unexportedFunc())
}

// ==================== init 函数 ====================

// 多个 init 函数按声明顺序执行
func init() {
	fmt.Println("[init 1] 第一个 init")
}

func init() {
	fmt.Println("[init 2] 第二个 init")
}

func demoInit() {
	fmt.Println("\n=== init 函数 ===")
	// init 已经在 main 之前执行了
}

// ==================== import 方式 ====================

// 演示不同的 import 方式（在代码中 import 只能在文件顶部）
func demoImports() {
	fmt.Println("\n=== import 方式 ===")

	// 标准导入：import "fmt"
	fmt.Println("标准导入: fmt.Println")

	// 别名导入：import f "fmt"  → f.Println()

	// 点导入（不推荐）：import . "fmt" → Println()  // 污染命名空间

	// 匿名导入：import _ "net/http/pprof"  // 仅执行 init
}

// ==================== internal 包 ====================

// Go 的 internal 包：放在 internal/ 目录下的包，
// 只能被其直接或间接父目录下的包导入。
// 例如：project/internal/config 只能被 project/ 下的包导入。

// ==================== go.mod 与模块 ====================

func demoModules() {
	fmt.Println("\n=== Go Modules ===")

	fmt.Println("go.mod 文件结构：")
	fmt.Println("  module GoLearning          // 模块路径")
	fmt.Println("  go 1.25.0                 // Go 版本")
	fmt.Println("  require (                 // 依赖")
	fmt.Println("    github.com/xxx v1.0.0")
	fmt.Println("  )")

	fmt.Println("\n常用命令：")
	fmt.Println("  go mod init <module-path>  // 初始化模块")
	fmt.Println("  go mod tidy                // 整理依赖")
	fmt.Println("  go get <package>           // 添加依赖")
	fmt.Println("  go mod vendor              // 创建 vendor 目录")
}

func main() {
	demoVisibility()
	demoInit()
	demoImports()
	demoModules()

	// 为了使用导入的包（避免编译错误）
	rand.Seed(time.Now().UnixNano())
	_ = rand.Int()
}
