package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	fmt.Println("=== Go 聊天客户端 ===")

	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Printf("无法连接到服务器: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("已连接到聊天服务器!")

	// 创建等待组
	var wg sync.WaitGroup

	// 启动接收消息的 goroutine
	wg.Add(1)
	go receiveMessages(conn, &wg)

	// 处理用户输入
	handleInput(conn)

	// 等待接收 goroutine 结束
	wg.Wait()
	fmt.Println("聊天客户端已退出")
}

func receiveMessages(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Print(message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("接收消息时出错: %v\n", err)
	}
}

func handleInput(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n可用命令:")
	fmt.Println("/list - 查看在线用户")
	fmt.Println("/quit 或 /exit - 退出聊天")
	fmt.Println("直接输入消息即可发送\n")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// 发送消息到服务器
		_, err := fmt.Fprintf(conn, "%s\n", input)
		if err != nil {
			fmt.Printf("发送消息失败: %v\n", err)
			break
		}

		// 检查退出命令
		if input == "/quit" || input == "/exit" {
			fmt.Println("正在退出...")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取输入时出错: %v\n", err)
	}
}
