// Package advanced 演示 Go 语言中的高级特性。
//
// 本文件涵盖：
//   - 构建标签（Build Tags）：条件编译
//   - go:embed：将静态文件嵌入到二进制
//   - unsafe：绕过 Go 类型系统
//   - 内存对齐
//   - cgo 概念
//   - go:generate：代码生成
//   - 编译器优化提示
package main

import (
	"fmt"
	"unsafe"
)

// ==================== 构建标签说明 ====================

// 构建标签（Build Tags）允许条件编译。
//
// 文件级别标签（Go 1.17+，推荐）：
//   //go:build linux
//   //go:build linux && amd64
//   //go:build !windows
//   //go:build linux || darwin
//
// 旧版标签（Go 1.16 及之前）：
//   // +build linux
//   // +build !windows,amd64
//
// 使用场景：
//   - 平台特定代码（_linux.go, _windows.go, _darwin.go）
//   - 架构特定代码（_amd64.go, _arm64.go）
//   - 测试专用代码（_test.go）
//   - 开发/生产代码区分

func demoBuildTags() {
	fmt.Println("=== 构建标签 ===")
	fmt.Println("构建标签允许为不同平台/架构编译不同代码")
	fmt.Println()
	fmt.Println("文件命名约定（自动识别）：")
	fmt.Println("  *_linux.go   - Linux 专属")
	fmt.Println("  *_windows.go - Windows 专属")
	fmt.Println("  *_darwin.go  - macOS 专属")
	fmt.Println("  *_amd64.go   - amd64 架构专属")
	fmt.Println("  *_arm64.go   - arm64 架构专属")
	fmt.Println("  *_test.go    - 测试文件")
	fmt.Println()
	fmt.Println("//go:build 指令示例：")
	fmt.Println("  //go:build linux")
	fmt.Println("  //go:build linux && amd64")
	fmt.Println("  //go:build !windows")
	fmt.Println("  //go:build linux || darwin")
}

// ==================== unsafe ====================

func demoUnsafe() {
	fmt.Println("\n=== unsafe 包 ===")

	// unsafe.Sizeof：获取变量占用的内存大小
	var (
		i  int
		f  float64
		s  string
		b  bool
		ch chan int
	)
	fmt.Printf("int 大小:     %d bytes\n", unsafe.Sizeof(i))
	fmt.Printf("float64 大小: %d bytes\n", unsafe.Sizeof(f))
	fmt.Printf("string 大小:  %d bytes (header: ptr+len)\n", unsafe.Sizeof(s))
	fmt.Printf("bool 大小:    %d bytes\n", unsafe.Sizeof(b))
	fmt.Printf("chan int 大小: %d bytes\n", unsafe.Sizeof(ch))

	// string 和 []byte 的内部结构
	// string header: {*byte, int}  -- 指针 + 长度
	// slice header:  {*byte, int, int} -- 指针 + 长度 + 容量

	// unsafe.Pointer：通用指针类型
	// 可用于不同类型指针的转换（慎用！）
	nums := []int{1, 2, 3, 4}
	ptr := unsafe.Pointer(&nums[0])
	fmt.Printf("\nunsafe.Pointer: %v\n", ptr)

	// 警告：unsafe.Pointer 的合法使用场景有限：
	// 1. 在 T1 和 T2 之间转换，前提是它们内存布局相同
	// 2. 在 unsafe.Pointer 和 uintptr 之间转换（用于指针运算）
	// 3. 调用 syscall.Syscall 时
}

// ==================== 内存对齐 ====================

// Unaligned 未优化的字段排列
type Unaligned struct {
	A bool  // 1 byte + 3 bytes padding
	B int64 // 8 bytes (需要 8 字节对齐)
	C bool  // 1 byte + 7 bytes padding
}

// Aligned 优化后的字段排列
type Aligned struct {
	B int64 // 8 bytes
	A bool  // 1 byte
	C bool  // 1 byte + 6 bytes padding
}

func demoMemoryAlignment() {
	fmt.Println("\n=== 内存对齐 ===")

	fmt.Printf("Unaligned: %d bytes (padding 浪费)\n", unsafe.Sizeof(Unaligned{}))
	fmt.Printf("Aligned:   %d bytes (更紧凑)\n", unsafe.Sizeof(Aligned{}))
	fmt.Println("提示: 将大小相近的字段放在一起可以减少 padding")

	// 更小的示例
	type Bad struct {
		flag    bool    // 1 + 7
		counter float64 // 8
	}
	type Good struct {
		counter float64 // 8
		flag    bool    // 1 + 7
	}
	fmt.Printf("Bad:  %d bytes\n", unsafe.Sizeof(Bad{}))
	fmt.Printf("Good: %d bytes\n", unsafe.Sizeof(Good{}))
}

// ==================== go:embed 说明 ====================

func demoEmbed() {
	fmt.Println("\n=== go:embed（需要 Go 1.16+） ===")
	fmt.Println("//go:embed 将文件嵌入到二进制中")
	fmt.Println()
	fmt.Println("用法示例：")
	fmt.Println("  import \"embed\"")
	fmt.Println("  //go:embed static/*")
	fmt.Println("  var staticFiles embed.FS")
	fmt.Println()
	fmt.Println("  //go:embed config.yaml")
	fmt.Println("  var configData []byte")
	fmt.Println()
	fmt.Println("  //go:embed templates/*.html")
	fmt.Println("  var templates embed.FS")
	fmt.Println()
	fmt.Println("支持：")
	fmt.Println("  - 嵌入为 string 或 []byte")
	fmt.Println("  - 嵌入为 embed.FS（文件系统）")
	fmt.Println("  - 支持通配符 */?")
	fmt.Println("  - 路径相对于源文件目录")
}

// ==================== go:generate 说明 ====================

func demoGenerate() {
	fmt.Println("=== go:generate ===")
	fmt.Println("//go:generate 用于自动生成代码")
	fmt.Println()
	fmt.Println("用法示例：")
	fmt.Println("  //go:generate stringer -type=Status")
	fmt.Println("  //go:generate mockgen -source=interface.go -destination=mock.go")
	fmt.Println("  //go:generate go run gen.go")
	fmt.Println()
	fmt.Println("运行：")
	fmt.Println("  go generate ./...")
}

// ==================== 编译器优化提示 ====================

func demoCompilerHints() {
	fmt.Println("\n=== 编译器优化提示 ===")

	// 提示1：尽量使用值类型而非指针（减少 GC 压力）
	// 提示2：字符串拼接用 strings.Builder（减少内存分配）
	// 提示3：slice 和 map 预先指定容量

	// 预先分配容量的对比
	const N = 10000

	// 差：频繁扩容
	bad := func() []int {
		var s []int
		for i := 0; i < N; i++ {
			s = append(s, i)
		}
		return s
	}

	// 好：预分配容量
	good := func() []int {
		s := make([]int, 0, N)
		for i := 0; i < N; i++ {
			s = append(s, i)
		}
		return s
	}

	_ = bad()
	_ = good()

	fmt.Println("预分配容量 vs 动态扩容")
	fmt.Println("提示: go test -bench=. -benchmem 查看内存分配")
}

// ==================== 优雅退出 ====================

func demoGracefulShutdown() {
	fmt.Println("\n=== 优雅退出模式 ===")
	fmt.Println("Go 程序建议监听系统信号实现优雅退出：")
	fmt.Println()
	fmt.Println("  quit := make(chan os.Signal, 1)")
	fmt.Println("  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)")
	fmt.Println("  <-quit")
	fmt.Println("  // 执行清理工作...")
	fmt.Println("  // server.Shutdown(ctx)")
	fmt.Println("  // db.Close()")
}

func main() {
	demoBuildTags()
	demoUnsafe()
	demoMemoryAlignment()
	demoEmbed()
	demoGenerate()
	demoCompilerHints()
	demoGracefulShutdown()
}
