// Package composite_types 演示 Go 语言中的结构体类型（struct）。
//
// 结构体的特点：
//   - 值类型：字段的集合，可以包含任意类型的字段
//   - 字段名首字母大写 = 导出（其他包可访问）；小写 = 包内私有
//   - 支持嵌入（embedding）：类似其他语言的「继承」，实际上是组合
//   - 支持标签（tag）：用于序列化、验证等元数据
//   - 零值：所有字段都是各自类型的零值
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// ==================== 结构体定义 ====================

// Person 定义一个人物结构体（导出类型、导出字段）
type Person struct {
	Name string `json:"name"` // 结构体标签：JSON 序列化时字段名为 "name"
	Age  int    `json:"age"`
}

// unexportedStruct 包内私有结构体（首字母小写）
type unexportedStruct struct {
	internalField string
}

// ==================== 结构体创建 ====================

func demoStructCreation() {
	fmt.Println("=== 结构体创建 ===")

	// 方式1：按顺序初始化所有字段（不推荐，字段顺序变化会出错）
	p1 := Person{"Alice", 30}
	fmt.Printf("按顺序: %+v\n", p1)

	// 方式2：按字段名初始化（推荐！）
	p2 := Person{
		Name: "Bob",
		Age:  25,
	}
	fmt.Printf("按字段名: %+v\n", p2)

	// 方式3：部分初始化（其余为零值）
	p3 := Person{Name: "Charlie"}
	fmt.Printf("部分初始化: %+v\n", p3)

	// 方式4：零值
	var p4 Person
	fmt.Printf("零值: %+v\n", p4)

	// 方式5：使用 new() 返回指针
	p5 := new(Person)
	p5.Name = "David"
	p5.Age = 35
	fmt.Printf("new() 指针: %+v\n", *p5)

	// 方式6：取地址 + 字面量（最常用的构造方式）
	p6 := &Person{
		Name: "Eve",
		Age:  28,
	}
	fmt.Printf("&字面量: %+v\n", *p6)
}

// ==================== 字段访问与修改 ====================

func demoFieldAccess() {
	fmt.Println("\n=== 字段访问与修改 ===")

	p := Person{Name: "Frank", Age: 40}

	// 访问字段
	fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)

	// 通过指针访问（自动解引用，不需要 *p）
	ptr := &p
	ptr.Name = "Grace" // 等同于 (*ptr).Name = "Grace"
	fmt.Printf("指针修改后: Name=%s\n", p.Name)

	// 修改字段
	p.Age = 41
	fmt.Printf("修改 Age 后: %+v\n", p)
}

// ==================== 结构体嵌入（组合） ====================

// Address 地址信息
type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

// Employee 雇员（嵌入 Person 和 Address）
type Employee struct {
	Person  // 嵌入 Person（匿名字段）
	Address // 嵌入 Address（匿名字段）
	Company string
	Salary  float64
}

func demoEmbedding() {
	fmt.Println("\n=== 结构体嵌入（组合） ===")

	e := Employee{
		Person:  Person{Name: "Henry", Age: 45},
		Address: Address{City: "Beijing", Country: "China"},
		Company: "Tech Corp",
		Salary:  100000,
	}

	// 直接访问嵌入类型的字段（提升后的字段）
	fmt.Printf("Name: %s, Age: %d, City: %s, Company: %s\n",
		e.Name, e.Age, e.City, e.Company)

	// 也可以显式访问（字段名即类型名）
	fmt.Printf("e.Person.Name: %s, e.Address.City: %s\n",
		e.Person.Name, e.Address.City)

	// 嵌入多个类型时，如果有字段名冲突，必须显式指定路径
}

// ==================== 结构体方法（值接收者 vs 指针接收者） ====================

// String 实现 fmt.Stringer 接口
func (p Person) String() string {
	return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

// Birthday 指针接收者方法：修改接收者的值
func (p *Person) Birthday() {
	p.Age++
}

// AgeInMonths 值接收者方法：不会修改接收者
func (p Person) AgeInMonths() int {
	return p.Age * 12
}

func demoMethods() {
	fmt.Println("\n=== 结构体方法 ===")

	p := Person{Name: "Ivy", Age: 30}

	// 值接收者方法
	fmt.Println(p.String())
	fmt.Printf("AgeInMonths: %d\n", p.AgeInMonths())

	// 指针接收者方法：修改原值
	p.Birthday()
	fmt.Printf("Birthday 后: %s\n", p.String())

	// Go 自动处理值/指针转换：
	// - 值类型可以调用指针接收者方法（编译器自动取地址）
	// - 指针类型可以调用值接收者方法（编译器自动解引用）
	ptr := &Person{Name: "Jack", Age: 20}
	fmt.Println(ptr.String()) // (*ptr).String() 自动解引用
}

// ==================== 结构体标签与 JSON ====================

// Product 商品（带多种标签）
type Product struct {
	ID        int       `json:"id" xml:"id"`
	Name      string    `json:"name" xml:"name" validate:"required"`
	Price     float64   `json:"price" xml:"price"`
	CreatedAt time.Time `json:"created_at"`
	Hidden    string    `json:"-"`                  // "-" 表示序列化时忽略
	Optional  string    `json:"optional,omitempty"` // omitempty: 零值时忽略
}

func demoTagsAndJSON() {
	fmt.Println("\n=== 结构体标签与 JSON ===")

	p := Product{
		ID:        1,
		Name:      "Laptop",
		Price:     9999.99,
		CreatedAt: time.Now(),
		Hidden:    "secret",
	}

	// 序列化为 JSON
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("JSON 序列化错误:", err)
		return
	}
	fmt.Printf("JSON: %s\n", string(jsonBytes))

	// 反序列化
	jsonStr := `{"id":2,"name":"Phone","price":5999.99}`
	var p2 Product
	json.Unmarshal([]byte(jsonStr), &p2)
	fmt.Printf("反序列化: %+v\n", p2)
}

// ==================== 空结构体 ====================

func demoEmptyStruct() {
	fmt.Println("\n=== 空结构体 ===")

	// struct{} 不占内存，常用于：
	// 1. 集合（set）
	// 2. 信号通道（signal channel）
	// 3. 仅用于方法的类型

	// 作为 set 的 value（不占用内存）
	set := make(map[string]struct{})
	set["key1"] = struct{}{}

	// 信号通道：chan struct{}
	done := make(chan struct{})
	go func() {
		// 做一些工作...
		close(done) // 发送完成信号
	}()
	<-done // 等待完成
	fmt.Println("空结构体信号通道已完成")

	fmt.Printf("struct{} 大小: %d byte\n", func() int {
		var s struct{}
		return int(unsafeGetSize(s))
	}())
}

// unsafeGetSize 演示 unsafe.Sizeof（简化版）
func unsafeGetSize(v interface{}) uintptr {
	// 实际使用: import "unsafe"; unsafe.Sizeof(v)
	return 0 // 实际 struct{} 大小为 0
}

func main() {
	demoStructCreation()
	demoFieldAccess()
	demoEmbedding()
	demoMethods()
	demoTagsAndJSON()
	demoEmptyStruct()
}
