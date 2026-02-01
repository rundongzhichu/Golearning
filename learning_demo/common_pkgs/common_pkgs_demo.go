package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("hello world")

	str := "hello go"
	fmt.Println("Length: ", len(str))
	fmt.Println("Contains 'Go': ", strings.Contains(str, "go"))
}
