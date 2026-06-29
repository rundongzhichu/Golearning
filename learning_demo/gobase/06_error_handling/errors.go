// Package error_handling 演示 Go 语言中的错误处理。
//
// Go 错误处理的特点：
//   - error 是一个内置接口，只有一个方法：Error() string
//   - 函数通常返回 (result, error)，调用方检查 error 是否为 nil
//   - 错误是值，不是异常（panic 才是异常）
//   - errors.New() 创建简单错误
//   - fmt.Errorf() 创建带格式的错误
//   - Go 1.13+：errors.Is() 判断错误链中是否包含特定错误
//   - Go 1.13+：errors.As() 从错误链中提取特定类型的错误
//   - Go 1.13+：fmt.Errorf + %w 创建包装错误（wrapping）
//   - Go 1.20+：errors.Join() 合并多个错误
package main

import (
	"errors"
	"fmt"
	"os"
)

// ==================== error 接口 ====================

// error 是内置接口：
// type error interface {
//     Error() string
// }

func demoErrorInterface() {
	fmt.Println("=== error 接口 ===")

	// errors.New() 创建错误
	err1 := errors.New("这是一个简单错误")
	fmt.Printf("errors.New: %v\n", err1)

	// fmt.Errorf() 创建带格式的错误
	name := "config.yaml"
	err2 := fmt.Errorf("无法打开文件 %s: %w", name, os.ErrNotExist)
	fmt.Printf("fmt.Errorf: %v\n", err2)

	// 通过 Error() 方法获取字符串
	fmt.Printf("err1.Error(): %s\n", err1.Error())
}

// ==================== 自定义错误类型 ====================

// ValidationError 验证错误
type ValidationError struct {
	Field string
	Value any
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("验证错误: 字段 '%s' 的值 '%v' %s", e.Field, e.Value, e.Msg)
}

// NotFoundError 未找到错误
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("未找到 %s: id=%s", e.Resource, e.ID)
}

func demoCustomError() {
	fmt.Println("\n=== 自定义错误 ===")

	// 创建不同类型的具体错误
	err1 := &ValidationError{Field: "email", Value: "not-an-email", Msg: "格式无效"}
	err2 := &NotFoundError{Resource: "用户", ID: "12345"}

	fmt.Printf("验证错误: %v\n", err1)
	fmt.Printf("未找到错误: %v\n", err2)
}

// ==================== 错误包装与解包（Go 1.13+） ====================

// 底层错误
var ErrDatabase = errors.New("数据库错误")
var ErrConnection = fmt.Errorf("连接失败: %w", ErrDatabase)

func getData() error {
	// 包装错误链：getData → ErrConnection → ErrDatabase
	return fmt.Errorf("getData 失败: %w", ErrConnection)
}

func demoErrorWrapping() {
	fmt.Println("\n=== 错误包装 (Go 1.13+) ===")

	err := getData()
	fmt.Printf("错误链: %v\n", err)

	// Unwrap：解包一层
	fmt.Printf("Unwrap: %v\n", errors.Unwrap(err))

	// errors.Is：检查错误链中是否包含特定错误（递归解包）
	fmt.Printf("Is ErrDatabase? %v\n", errors.Is(err, ErrDatabase))
	fmt.Printf("Is ErrConnection? %v\n", errors.Is(err, ErrConnection))
	fmt.Printf("Is os.ErrNotExist? %v\n", errors.Is(err, os.ErrNotExist))

	// errors.As：从错误链中提取特定类型的错误
	var validationErr *ValidationError
	testErr := &ValidationError{Field: "age", Value: -1, Msg: "不能为负数"}
	fmt.Printf("As ValidationError: %v\n", errors.As(testErr, &validationErr))
}

// ==================== errors.Join（Go 1.20+） ====================

func demoErrorJoin() {
	fmt.Println("\n=== errors.Join (Go 1.20+) ===")

	err1 := errors.New("错误1")
	err2 := errors.New("错误2")
	err3 := errors.New("错误3")

	// 合并多个错误
	combined := errors.Join(err1, err2, err3)
	fmt.Printf("合并后: %v\n", combined)

	// errors.Is 仍然有效
	fmt.Printf("Is err1? %v\n", errors.Is(combined, err1))
	fmt.Printf("Is err2? %v\n", errors.Is(combined, err2))

	// 合并 nil 错误：nil 会被过滤掉
	clean := errors.Join(err1, nil, err3)
	fmt.Printf("合并(含 nil): %v\n", clean)
}

// ==================== 错误处理常见模式 ====================

// 模式1：简单传递（不包装）
func pattern1() error {
	_, err := os.Open("/nonexistent")
	if err != nil {
		return err // 直接返回
	}
	return nil
}

// 模式2：包装错误（添加上下文）
func pattern2() error {
	_, err := os.Open("/nonexistent")
	if err != nil {
		return fmt.Errorf("pattern2: 打开文件失败: %w", err) // 包装
	}
	return nil
}

// 模式3：检查并处理特定错误
func pattern3() {
	_, err := os.Open("/nonexistent")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("  文件不存在，使用默认配置")
		} else {
			fmt.Printf("  其他错误: %v\n", err)
		}
	}
}

// 模式4：defer 中的错误处理
func pattern4() (err error) {
	f, err := os.Open("/nonexistent")
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if closeErr != nil && err == nil {
			err = closeErr // 只在原来没有错误时返回 Close 的错误
		}
	}()
	return nil
}

func demoErrorPatterns() {
	fmt.Println("\n=== 错误处理模式 ===")

	fmt.Println("模式1 (直接返回):")
	if err := pattern1(); err != nil {
		fmt.Printf("  err: %v\n", err)
	}

	fmt.Println("模式2 (包装错误):")
	if err := pattern2(); err != nil {
		fmt.Printf("  err: %v\n", err)
	}

	fmt.Println("模式3 (检查特定错误):")
	pattern3()

	fmt.Println("模式4 (defer 错误处理):")
	_ = pattern4()
}

func main() {
	demoErrorInterface()
	demoCustomError()
	demoErrorWrapping()
	demoErrorJoin()
	demoErrorPatterns()
}
