// Package composite_types 演示 Go 语言中的数组类型。
//
// 数组的特点：
//   - 长度固定，是类型的一部分：[3]int 和 [5]int 是不同的类型
//   - 值类型：赋值和传参会复制整个数组
//   - 零值：每个元素都是对应类型的零值
//   - 数组的长度可以用 len() 获取
//   - 实际开发中，切片（slice）更常用
package main

import "fmt"

func demoArrayDeclaration() {
	fmt.Println("=== 数组声明 ===")

	// 声明方式1：指定长度，使用零值
	var arr1 [5]int
	fmt.Printf("var [5]int 零值: %v\n", arr1)

	// 声明方式2：声明并初始化
	var arr2 [3]string = [3]string{"Go", "Python", "Rust"}
	fmt.Printf("var [3]string: %v\n", arr2)

	// 声明方式3：短变量声明 + 字面量
	arr3 := [4]float64{1.1, 2.2, 3.3, 4.4}
	fmt.Printf("短声明 + 字面量: %v\n", arr3)

	// 声明方式4：使用 ... 让编译器推断长度
	arr4 := [...]int{1, 2, 3, 4, 5, 6}
	fmt.Printf("编译器推断长度: %v (len=%d)\n", arr4, len(arr4))

	// 声明方式5：指定索引初始化（其余为零值）
	arr5 := [5]int{0: 10, 4: 50}
	fmt.Printf("按索引初始化:   %v\n", arr5)
}

func demoArrayAccess() {
	fmt.Println("\n=== 数组访问与遍历 ===")

	arr := [5]int{10, 20, 30, 40, 50}

	// 下标访问（0-based）
	fmt.Printf("arr[0]=%d, arr[4]=%d\n", arr[0], arr[4])

	// 修改元素
	arr[2] = 100
	fmt.Printf("修改后: %v\n", arr)

	// 长度
	fmt.Printf("len(arr)=%d\n", len(arr))

	// 遍历方式1：for i
	fmt.Print("for i: ")
	for i := 0; i < len(arr); i++ {
		fmt.Printf("%d ", arr[i])
	}
	fmt.Println()

	// 遍历方式2：for range（推荐）
	fmt.Print("for range: ")
	for i, v := range arr {
		fmt.Printf("[%d]=%d ", i, v)
	}
	fmt.Println()

	// 只取值，忽略索引
	for _, v := range arr {
		_ = v
	}
}

func demoArrayValueSemantics() {
	fmt.Println("\n=== 数组的值语义 ===")

	// 数组是值类型：赋值会复制整个数组
	arr1 := [3]int{1, 2, 3}
	arr2 := arr1  // 完整复制
	arr2[0] = 999 // 修改 arr2 不影响 arr1

	fmt.Printf("arr1: %v (未受影响)\n", arr1)
	fmt.Printf("arr2: %v (已修改)\n", arr2)

	// 函数传参同样会复制（大数组时注意性能）
	// 如需引用传递，使用 *[N]T 或切片
	modifyArray := func(arr [3]int) {
		arr[0] = 777 // 只修改副本
	}
	modifyArray(arr1)
	fmt.Printf("函数传参后 arr1: %v (未受影响)\n", arr1)

	// 使用指针传递以避免拷贝
	modifyArrayPtr := func(arr *[3]int) {
		arr[0] = 777 // 通过指针修改原数组
	}
	modifyArrayPtr(&arr1)
	fmt.Printf("指针传参后 arr1: %v (已修改)\n", arr1)
}

func demoMultiDimensionArray() {
	fmt.Println("\n=== 多维数组 ===")

	// 2x3 矩阵
	var matrix [2][3]int
	matrix[0] = [3]int{1, 2, 3}
	matrix[1] = [3]int{4, 5, 6}

	fmt.Printf("二维数组: %v\n", matrix)
	fmt.Printf("matrix[1][2] = %d\n", matrix[1][2])

	// 遍历二维数组
	fmt.Println("遍历二维数组:")
	for i, row := range matrix {
		for j, val := range row {
			fmt.Printf("  [%d][%d] = %d\n", i, j, val)
		}
	}

	// 声明并初始化二维数组
	grid := [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Printf("grid: %v\n", grid)
}

func main() {
	demoArrayDeclaration()
	demoArrayAccess()
	demoArrayValueSemantics()
	demoMultiDimensionArray()
}
