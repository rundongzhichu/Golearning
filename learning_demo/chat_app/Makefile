# Go 聊天程序 Makefile
# 这段代码定义了Makefile中的伪目标（phony targets），包括 `build`、`run-server`、`run-client` 和 `clean`
# 这些目标不对应实际文件，用于执行特定任务，如构建项目、运行服务端/客户端或清理生成文件。
.PHONY: build run-server run-client clean

# 编译所有程序
build:
	go build -o chat-server server.go
	go build -o chat-client client.go
	@echo "编译完成!"

# 运行服务器
run-server:
	go run server.go

# 运行客户端
run-client:
	go run client.go

# 清理编译文件
clean:
	rm -f chat-server chat-client
	@echo "清理完成!"

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make build     - 编译服务器和客户端"
	@echo "  make run-server - 运行聊天服务器"
	@echo "  make run-client - 运行聊天客户端"
	@echo "  make clean     - 清理编译文件"
	@echo "  make help      - 显示此帮助信息"