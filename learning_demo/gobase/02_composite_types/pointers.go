// Package composite_types 演示 Go 语言中的指针类型。
//
// 指针的特点：
//   - 存储另一个变量的内存地址
//   - 零值为 nil
//   - 使用 & 获取变量的地址，使用 * 解引用
//   - Go 没有指针运算（除非使用 unsafe 包）
//   - 指针用于在函数间共享数据，避免大对象拷贝
//   - new() 函数分配内存并返回指针
package main

import "fmt"

// ==================== 指针基础 ====================

func demoPointerBasics() {
	fmt.Println("=== 指针基础 ===")

	// 声明指针（零值为 nil）
	var ptr *int
	fmt.Printf("nil 指针: %v (地址: %p)\n", ptr, ptr)

	// 取地址 &
	x := 42
	ptr = &x
	fmt.Printf("x=%d, &x=%p, ptr=%p, *ptr=%d\n", x, &x, ptr, *ptr)

	// 通过指针修改值
	*ptr = 100
	fmt.Printf("修改 *ptr=100 后, x=%d\n", x)

	// new() 分配零值并返回指针
	y := new(int) // y 是 *int，指向值为 0 的 int
	fmt.Printf("new(int): 值=%d, 地址=%p\n", *y, y)
	*y = 200
	fmt.Printf("赋值后: *y=%d\n", *y)
}

// ==================== 指针与函数 ====================

// modifyValue 值传递：不会影响原变量
func modifyValue(n int) {
	n = 999
}

// modifyPointer 指针传递：会影响原变量
func modifyPointer(n *int) {
	*n = 999
}

func demoPointerWithFunctions() {
	fmt.Println("\n=== 指针与函数 ===")

	x := 100

	// 值传递：原变量不变
	modifyValue(x)
	fmt.Printf("值传递后: x=%d\n", x)

	// 指针传递：原变量被修改
	modifyPointer(&x)
	fmt.Printf("指针传递后: x=%d\n", x)
}

// ==================== 指针与结构体 ====================

type Config struct {
	Host string
	Port int
}

// updateConfig 通过指针修改结构体
func updateConfig(cfg *Config) {
	cfg.Host = "localhost"
	cfg.Port = 8080
}

func demoPointerWithStructs() {
	fmt.Println("\n=== 指针与结构体 ===")

	cfg := Config{Host: "0.0.0.0", Port: 9090}
	fmt.Printf("修改前: %+v\n", cfg)

	updateConfig(&cfg)
	fmt.Printf("修改后: %+v\n", cfg)

	// 结构体指针访问字段时自动解引用
	ptr := &cfg
	fmt.Printf("ptr.Host=%s (等同于 (*ptr).Host)\n", ptr.Host)
}

// ==================== 指针与切片/map ====================

func demoPointerWithReferenceTypes() {
	fmt.Println("\n=== 指针与引用类型 ===")

	// slice、map、channel 本身是引用类型，通常不需要指针传递

	// slice 示例：函数内修改会影响原 slice（共享底层数组）
	modifySlice := func(s []int) {
		s[0] = 999
	}
	nums := []int{1, 2, 3}
	modifySlice(nums)
	fmt.Printf("slice 修改后: %v\n", nums)

	// 但如果要修改 slice 本身（如 append 后重新赋值），需要用指针
	appendSlice := func(s *[]int) {
		*s = append(*s, 4, 5)
	}
	appendSlice(&nums)
	fmt.Printf("指针 append 后: %v\n", nums)

	// map 示例：不需要指针
	modifyMap := func(m map[string]int) {
		m["new_key"] = 100
	}
	m := map[string]int{"old_key": 1}
	modifyMap(m)
	fmt.Printf("map 修改后: %v\n", m)
}

// ==================== 指针数组 vs 数组指针 ====================

func demoPointerArrays() {
	fmt.Println("\n=== 指针数组 vs 数组指针 ===")

	// 指针数组：数组的元素是指针
	a, b, c := 1, 2, 3
	ptrArray := [3]*int{&a, &b, &c}
	fmt.Printf("指针数组: [0]=%d, [1]=%d, [2]=%d\n",
		*ptrArray[0], *ptrArray[1], *ptrArray[2])

	// 数组指针：指向数组的指针
	arr := [3]int{10, 20, 30}
	arrPtr := &arr
	arrPtr[0] = 100 // 自动解引用
	fmt.Printf("数组指针修改后: arr=%v\n", arr)
}

// ==================== 二级指针 ====================

func demoDoublePointer() {
	fmt.Println("\n=== 二级指针 ===")

	x := 10
	p1 := &x  // *int
	p2 := &p1 // **int

	fmt.Printf("x=%d, *p1=%d, **p2=%d\n", x, *p1, **p2)

	**p2 = 200
	fmt.Printf("通过二级指针修改后: x=%d\n", x)
}

// ==================== 指针的常见使用模式 ====================

type LargeStruct struct {
	Data [1024]byte // 模拟大结构体
}

// processByValue 值传递会复制整个大结构体（性能差）
func processByValue(ls LargeStruct) {
	_ = ls.Data[0]
}

// processByPointer 指针传递只复制 8 字节（64 位系统）
func processByPointer(ls *LargeStruct) {
	_ = ls.Data[0]
}

func demoCommonPatterns() {
	fmt.Println("\n=== 指针常用模式 ===")

	// 模式1：从函数返回指针（Go 支持返回局部变量的指针）
	newPerson := func(name string, age int) *Person1 {
		// p 在栈上分配，但返回指针后自动逃逸到堆
		return &Person1{name, age}
	}
	p := newPerson("Alice", 30)
	fmt.Printf("返回局部指针: %+v\n", p)

	// 模式2：nil 指针检查
	var ptr *int
	if ptr == nil {
		fmt.Println("nil 检查通过：指针为零值")
	}

	// 模式3：可选参数使用指针类型
	// 当需要区分 "未设置" 和 "零值" 时使用指针
	config := struct {
		Timeout *int // nil 表示未设置，非 nil 表示显式设置了值
	}{}
	if config.Timeout == nil {
		fmt.Println("Timeout 未设置")
	}
	t := 30
	config.Timeout = &t
	fmt.Printf("Timeout 设置为 %d\n", *config.Timeout)
}

// Person1 简单的人员结构体
type Person1 struct {
	Name string
	Age  int
}

func main() {
	demoPointerBasics()
	demoPointerWithFunctions()
	demoPointerWithStructs()
	demoPointerWithReferenceTypes()
	demoPointerArrays()
	demoDoublePointer()
	demoCommonPatterns()
}
