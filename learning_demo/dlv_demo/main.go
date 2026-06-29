package main

import "fmt"

/*
dlv调试：
dlv debug main.go

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
