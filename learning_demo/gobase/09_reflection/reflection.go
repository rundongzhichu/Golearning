// Package reflection 演示 Go 语言中的反射（reflect 包）。
//
// 反射的特点：
//   - 允许程序在运行时检查、修改自身的行为和类型
//   - reflect.Type：类型信息（通过 reflect.TypeOf() 获取）
//   - reflect.Value：值信息（通过 reflect.ValueOf() 获取）
//   - 可以通过反射修改值（前提是可寻址、可设置）
//   - 结构体标签（tag）通过反射读取
//   - 反射性能较低，不应滥用；优先使用泛型
//
// 反射三定律（出自 Go 官方博客）：
//  1. 反射可以从接口值得到反射对象
//  2. 反射可以从反射对象得到接口值
//  3. 要修改反射对象，其值必须是可设置的（settable）
package main

import (
	"fmt"
	"reflect"
	"strings"
)

// ==================== 类型与值 ====================

type Person struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age"`
}

// Greet Person 的方法
func (p Person) Greet(greeting string) string {
	return fmt.Sprintf("%s, I'm %s, %d years old", greeting, p.Name, p.Age)
}

func demoTypeAndValue() {
	fmt.Println("=== reflect.Type 和 reflect.Value ===")

	p := Person{Name: "Alice", Age: 30}

	// 获取类型
	t := reflect.TypeOf(p)
	fmt.Printf("Type: %v (Kind: %v)\n", t, t.Kind())

	// 获取值
	v := reflect.ValueOf(p)
	fmt.Printf("Value: %v (Type: %v)\n", v, v.Type())

	// 从反射对象获取接口值
	i := v.Interface()
	fmt.Printf("Interface: %v\n", i)
}

// ==================== 查看结构体字段 ====================

func demoStructFields() {
	fmt.Println("\n=== 结构体字段反射 ===")

	p := Person{Name: "Bob", Age: 25}
	t := reflect.TypeOf(p)

	// 遍历字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("  字段 %d: %s (类型: %v, 标签: %q)\n",
			i, field.Name, field.Type, field.Tag)

		// 提取标签
		jsonTag := field.Tag.Get("json")
		validateTag := field.Tag.Get("validate")
		fmt.Printf("    json: %s, validate: %s\n", jsonTag, validateTag)
	}
}

// ==================== 查看和调用方法 ====================

func demoMethods() {
	fmt.Println("\n=== 方法反射 ===")

	p := Person{Name: "Charlie", Age: 28}
	t := reflect.TypeOf(p)

	// 遍历方法
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("  方法 %d: %s\n", i, method.Name)
		fmt.Printf("    类型: %v\n", method.Type)
		fmt.Printf("    参数数: %d\n", method.Type.NumIn())
		fmt.Printf("    返回值数: %d\n", method.Type.NumOut())
	}

	// 调用方法
	v := reflect.ValueOf(p)
	method := v.MethodByName("Greet")
	if method.IsValid() {
		args := []reflect.Value{reflect.ValueOf("Hi")}
		results := method.Call(args)
		fmt.Printf("\n调用 Greet: %v\n", results[0].String())
	}
}

// ==================== 修改值 ====================

func demoSetValue() {
	fmt.Println("\n=== 修改值（Settable） ===")

	x := 42

	// 错误：v 不可设置（因为传的是副本）
	v := reflect.ValueOf(x)
	fmt.Printf("CanSet: %v (因为传的是值副本)\n", v.CanSet())

	// 正确：传递指针
	v = reflect.ValueOf(&x).Elem()
	fmt.Printf("CanSet: %v\n", v.CanSet())

	// 修改值
	if v.CanSet() {
		v.SetInt(100)
	}
	fmt.Printf("x = %d\n", x)

	// 修改结构体字段
	p := Person{Name: "David", Age: 30}
	pv := reflect.ValueOf(&p).Elem()
	nameField := pv.FieldByName("Name")
	if nameField.IsValid() && nameField.CanSet() {
		nameField.SetString("Eve")
	}
	fmt.Printf("修改后: %+v\n", p)
}

// ==================== 实用场景 ====================

// StructToMap 将结构体转为 map[string]any（通过反射）
func StructToMap(s any) map[string]any {
	result := make(map[string]any)
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 获取 json 标签作为 key
		key := field.Tag.Get("json")
		if key == "" || key == "-" {
			key = strings.ToLower(field.Name)
		} else {
			// 去掉 omitempty 等选项
			if comma := strings.Index(key, ","); comma != -1 {
				key = key[:comma]
			}
		}
		result[key] = value.Interface()
	}
	return result
}

func demoStructToMap() {
	fmt.Println("\n=== 反射实用场景 ===")

	p := Person{Name: "Alice", Age: 30}
	m := StructToMap(p)
	fmt.Printf("StructToMap: %v\n", m)
}

// ==================== reflect.DeepEqual ====================

func demoDeepEqual() {
	fmt.Println("\n=== reflect.DeepEqual ===")

	// 比较结构体
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}
	p3 := Person{Name: "Alice", Age: 31}

	fmt.Printf("p1 == p2: %v\n", reflect.DeepEqual(p1, p2))
	fmt.Printf("p1 == p3: %v\n", reflect.DeepEqual(p1, p3))

	// 比较切片
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	fmt.Printf("s1 == s2: %v\n", reflect.DeepEqual(s1, s2))

	// 比较 map
	m1 := map[string]int{"a": 1}
	m2 := map[string]int{"a": 1}
	fmt.Printf("m1 == m2: %v\n", reflect.DeepEqual(m1, m2))
}

// ==================== Kind 类型 ====================

func demoKind() {
	fmt.Println("\n=== reflect.Kind ===")

	tests := []any{
		42,
		"hello",
		true,
		[]int{1, 2, 3},
		map[string]int{"a": 1},
		struct{ Name string }{"test"},
		func() {},
		make(chan int),
	}

	for _, v := range tests {
		fmt.Printf("  %T → Kind: %s\n", v, reflect.TypeOf(v).Kind())
	}
}

func main() {
	demoTypeAndValue()
	demoStructFields()
	demoMethods()
	demoSetValue()
	demoStructToMap()
	demoDeepEqual()
	demoKind()
}
