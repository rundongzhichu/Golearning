package main

import "fmt"

/*
 1， go build -gcflags="-N -l" -o test main.go

gdb调试：
gdb ./test # 启动调试
gdb --args ./test arg1 arg2 # 指定参数启动调试

目前本地macOs 安装的gdb不支持本地调试，需要配置证书，
*/

func add(a, b int) int {
	sum := 0
	sum = a + b
	return sum
}
func main() {
	sum := add(10, 20)
	fmt.Println(sum)
}
