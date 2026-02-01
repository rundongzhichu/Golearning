package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn net.Conn
	name string
}

// 全局变量用于管理所有连接的客户端
var (
	// clients 存储所有活跃的客户端连接，键为网络连接，值为客户端信息
	clients = make(map[net.Conn]*Client)
	// clientsMu 保护 clients map 的读写操作，确保并发安全
	clientsMu sync.RWMutex
)

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("聊天服务器已启动，监听端口 8081...")
	fmt.Println("等待客户端连接...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败: %v\n", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// 获取客户端地址
	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("新客户端连接: %s\n", clientAddr)

	// 创建客户端对象
	client := &Client{
		conn: conn,
		name: clientAddr,
	}

	// 添加到客户端列表
	clientsMu.Lock()
	clients[conn] = client
	clientsMu.Unlock()

	// 发送欢迎消息
	welcomeMsg := fmt.Sprintf("欢迎来到聊天室! 您的ID是: %s\n请输入您的昵称: ", clientAddr)
	conn.Write([]byte(welcomeMsg))

	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// 读取昵称
	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取昵称失败: %v\n", err)
		return
	}

	// 清理昵称
	name = strings.TrimSpace(name)
	if name != "" {
		client.name = name
	}

	// 取消读取超时
	conn.SetReadDeadline(time.Time{})

	fmt.Printf("客户端 %s 设置昵称为: %s\n", clientAddr, client.name)

	// 通知其他用户
	broadcastMessage(fmt.Sprintf("*** %s 加入了聊天室 ***\n", client.name), conn)

	// 处理消息循环
	handleMessages(client)

	// 客户端断开连接
	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()

	broadcastMessage(fmt.Sprintf("*** %s 离开了聊天室 ***\n", client.name), nil)
	fmt.Printf("客户端断开连接: %s (%s)\n", client.name, clientAddr)
}

func handleMessages(client *Client) {
	reader := bufio.NewReader(client.conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// 处理特殊命令
		if message == "/quit" || message == "/exit" {
			client.conn.Write([]byte("再见!\n"))
			break
		}

		if message == "/list" {
			listClients(client)
			continue
		}

		// 广播消息
		timestamp := time.Now().Format("15:04:05")
		formattedMessage := fmt.Sprintf("[%s] %s: %s\n", timestamp, client.name, message)
		broadcastMessage(formattedMessage, client.conn)
	}
}

func broadcastMessage(message string, excludeConn net.Conn) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	for conn, client := range clients {
		if conn != excludeConn {
			_, err := conn.Write([]byte(message))
			if err != nil {
				fmt.Printf("向客户端 %s 发送消息失败: %v\n", client.name, err)
			}
		}
	}
}

func listClients(client *Client) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	client.conn.Write([]byte("=== 在线用户列表 ===\n"))
	for _, c := range clients {
		status := "在线"
		client.conn.Write([]byte(fmt.Sprintf("- %s (%s) %s\n", c.name, c.conn.RemoteAddr().String(), status)))
	}
	client.conn.Write([]byte("==================\n"))
}
