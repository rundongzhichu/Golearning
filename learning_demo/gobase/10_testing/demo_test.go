// Package testing 演示 Go 语言中的测试。
//
// Go 测试的特点：
//   - 测试文件以 _test.go 结尾
//   - 测试函数以 Test 开头，接收 *testing.T
//   - 基准测试函数以 Benchmark 开头，接收 *testing.B
//   - 示例函数以 Example 开头
//   - 使用 go test 命令运行测试
//   - 表驱动测试（table-driven tests）是 Go 社区推荐的测试模式
//   - testing.T 提供 Log、Error、Fatal、Run（子测试）等方法
//   - t.Parallel() 标记测试可并行执行
//   - testing.B 用于基准测试，b.N 控制迭代次数
//   - testing.F 用于模糊测试（Fuzzing），Go 1.18+
package main

import (
	"fmt"
	"os"
	"testing"
)

// ==================== 业务代码 ====================

// Add 加法（被测试的函数）
func Add(a, b int) int {
	return a + b
}

// Divide 除法（返回 error 的函数）
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}
	return a / b, nil
}

// ==================== 基础测试 ====================

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5
	if result != expected {
		t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
	}
}

// ==================== 表驱动测试 ====================

func TestAdd_TableDriven(t *testing.T) {
	// 定义测试用例表
	tests := []struct {
		name     string // 测试名称
		a, b     int    // 输入
		expected int    // 期望输出
	}{
		{"正数相加", 2, 3, 5},
		{"零相加", 0, 0, 0},
		{"负数相加", -1, -1, -2},
		{"混合相加", -1, 1, 0},
		{"大数相加", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; expected %d",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// ==================== 错误测试 ====================

func TestDivide(t *testing.T) {
	// 正常情况
	result, err := Divide(10, 2)
	if err != nil {
		t.Errorf("Divide(10,2) 不应该出错: %v", err)
	}
	if result != 5 {
		t.Errorf("Divide(10,2) = %d; expected 5", result)
	}

	// 除零情况
	_, err = Divide(10, 0)
	if err == nil {
		t.Error("Divide(10,0) 应该返回错误")
	}
}

// ==================== 子测试 ====================

func TestDivide_Cases(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int
		expected  int
		expectErr bool
	}{
		{"正常除法", 10, 2, 5, false},
		{"除零", 10, 0, 0, true},
		{"负数除法", 10, -2, -5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			if tt.expectErr {
				if err == nil {
					t.Error("期望出错但没有错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望出错但出错: %v", err)
				}
				if result != tt.expected {
					t.Errorf("=%d, 期望 %d", result, tt.expected)
				}
			}
		})
	}
}

// ==================== 并行测试 ====================

func TestParallel(t *testing.T) {
	t.Parallel() // 标记为可并行

	// 多个并行子测试
	for i := 0; i < 3; i++ {
		i := i // 捕获循环变量
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()
			result := Add(i, 1)
			if result != i+1 {
				t.Error("结果不对")
			}
		})
	}
}

// ==================== 基准测试 ====================

// BenchmarkAdd 基准测试：go test -bench=.
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(2, 3)
	}
}

// ==================== 示例测试 ====================

// ExampleAdd 示例函数：执行并验证输出
// Output 注释用于验证标准输出
func ExampleAdd() {
	fmt.Println(Add(1, 2))
	// Output: 3
}

// ==================== TestMain ====================

// TestMain 测试入口：在测试前后执行 setup/teardown
func TestMain(m *testing.M) {
	fmt.Println("=== 测试开始 ===")
	// setup code...
	code := m.Run() // 执行所有测试
	// teardown code...
	fmt.Println("=== 测试结束 ===")
	os.Exit(code)
}

// ==================== 辅助函数 ====================

// assertEqual 测试辅助函数（使用 t.Helper）
func assertEqual(t *testing.T, got, expected int) {
	t.Helper() // 标记为辅助函数，错误报告会指向调用方
	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}

func TestWithHelper(t *testing.T) {
	result := Add(3, 7)
	assertEqual(t, result, 10)
}
