package main

import (
	"fmt"
	"os"
)

// go 函数
func greet(name string) {
	fmt.Println("Hello, " + name)
}

// go结构体
type Person struct {
	name string
	age  int
}

func print_numbers(ch chan int) {
	for i := 1; i < 5; i++ {
		ch <- i
	}
	close(ch)
}

// 包内类型简洁命名（外部通过包名访问）
type Counter struct {
	count int
}

// 导出方法首字母大写
func (c *Counter) Increment() {
	c.count++
}

func (c Counter) Value() int {
	return c.count
}

func readFile(filename string) {
	fmt.Println("Opening", filename)
	defer fmt.Println("Closing", filename) // 最后执行
	// 模拟读取
	fmt.Println("Reading content...")
}

// 定义小接口
type Stringer interface {
	String() string
}

// 任何类型只要实现了String()就满足Stringer接口
func (p Person) String() string {
	return fmt.Sprintf("%s (%d years)", p.name, p.age)
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

var config struct {
	Port string
}

func init() {
	// 初始化配置
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}
	fmt.Println("Config initialized:", config.Port)
}

func main() {
	// 基础变量声明
	var name string = "Alice" // 显式声明
	age := 25                 // 简短声明
	fmt.Println(name, age)

	// 函数调用
	greet(name)

	// 数组切片
	var arr [3]int = [3]int{1, 2, 3} // 数组是固定大小的
	fmt.Println(arr)

	slice := []int{1, 2, 3, 4} // 切片是可变长度
	// 使用内置 append
	slice = append(slice, 4, 5)
	// 追加另一个切片（注意 ...）
	more := []int{6, 7}
	slice = append(slice, more...)
	fmt.Println(slice)
	// if for switch
	number := 3
	if number%2 == 0 {
		fmt.Println("Even")
	} else {
		fmt.Println("Odd")
	}

	for i := 1; i <= 5; i++ {
		fmt.Println(i)
	}

	switch number {
	case 1:
		fmt.Println("One")
	case 2:
		fmt.Println("Two")
	default:
		fmt.Println("Other")
	}

	// go 结构体
	alice := Person{name: "Alice", age: 27}
	fmt.Println(alice.name, alice.age)

	// go 并发编程
	/**
	Go的并发模型基于goroutines和channels。
	goroutine：是Go中的轻量级线程，可以通过go关键字启动。
	channel：用于goroutine之间通信。
	*/
	ch := make(chan int)
	go print_numbers(ch)
	for i := range ch {
		fmt.Println(i)
	}

	// 结构体方法
	var c Counter
	c.Increment()
	fmt.Println("Count:", c.Value()) // Count: 1

	// if with initializer  if错误处理
	if file, err := os.Open("nonexistent.txt"); err != nil {
		fmt.Println("Error:", err) // 正确处理错误
	} else {
		defer file.Close()
	}

	// for-range 遍历 map
	m := map[string]int{"a": 1, "b": 2}
	for k, v := range m {
		fmt.Printf("%s: %d\n", k, v)
	}

	// defer资源处理
	// demo_defer.go
	readFile("data.txt")
	// 输出：
	// Opening data.txt
	// Reading content...
	// Closing data.txt

	// new  vs  make
	// new: 返回指针，内容为零值
	/*
		new([]int)：分配一个指向空切片的指针，初始值为nil。
		make([]int, 3, 10)：创建长度为3、容量为10的切片，并初始化为零值。
		make(map[string]int)：创建一个空映射。
	*/
	p := new([]int)
	fmt.Println("*p is nil?", *p == nil) // true
	fmt.Println("*p is nil?", p)         // true

	// make: 初始化 slice/map/channel
	s := make([]int, 3, 10)
	fmt.Println("slice:", s) // [100 0 0]
	s[0] = 100
	fmt.Println("slice:", s) // [100 0 0]

	m1 := make(map[string]int)
	fmt.Println("map:", m1) // map[key:42]
	m1["key"] = 42
	fmt.Println("map:", m1) // map[key:42]

	// Stringer接口实现
	p1 := Person{"Alice", 30}
	var s1 Stringer = p1 // 隐式实现
	fmt.Println(s1)      // Alice (30 years)

	// _ 忽略变量
	// 只关心错误
	if _, err := divide(10, 0); err != nil {
		fmt.Println("Error:", err)
	}

	// 忽略 map 中不存在的 key
	m2 := map[string]int{"x": 1}
	if _, exists := m2["y"]; !exists {
		fmt.Println("Key 'y' not found", exists)
	}

	// 初始化函数
	fmt.Println("Server starting on port", config.Port)
}
