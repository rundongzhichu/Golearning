# Golearning
这是一个用来学习go语言的项目

# 🚀 Go 语言基础语法与工程化开发实战教程

> 适合有编程经验的开发者快速掌握 Go 核心语法 + 现代工程实践  
> **版本**：Go 1.22+ | **目标**：写出可维护、可测试、可部署的生产级 Go 代码

---

## 第一部分：Go 基础语法（30 分钟速成）

### 1. Hello World & 基本结构
```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```
- **包**（package）：每个 Go 文件必须属于一个包，可执行程序入口为 `main` 包
- **导入**（import）：标准库无需路径，第三方库用模块路径（如 `github.com/gin-gonic/gin`）

---

### 2. 变量与类型
```go
// 声明方式
var name string = "Go"
age := 15 // := 自动推导类型（仅函数内可用）

// 基础类型
bool, string
int, int8, int16, int32, int64
uint, uintptr
float32, float64
complex64, complex128

// 特殊类型
byte // uint8 别名
rune // int32 别名（表示 Unicode 码点）
```

---

### 3. 控制结构
```go
// if (无需括号)
if x := 10; x > 5 {
    fmt.Println(x)
}

// for (Go 只有 for 循环)
for i := 0; i < 10; i++ { }
for key, value := range mapVar { } // 遍历 map/slice

// switch (无需 break)
switch os := runtime.GOOS; os {
case "darwin":
    fmt.Println("Mac")
default:
    fmt.Println(os)
}
```

---

### 4. 函数与错误处理
```go
// 多返回值（惯用法：(result, error)）
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// 调用
result, err := divide(10, 2)
if err != nil {
    log.Fatal(err)
}
```
> ✅ **Go 哲学**：显式错误处理，拒绝异常（panic 仅用于不可恢复错误）

---

### 5. 结构体与方法
```go
// 结构体定义
type User struct {
    ID   int    `json:"id"`     // 标签（tag）
    Name string `json:"name"`
}

// 方法（接收者）
func (u User) Greet() string {
    return "Hello, " + u.Name
}

// 使用
user := User{ID: 1, Name: "Alice"}
fmt.Println(user.Greet()) // Hello, Alice
```

---

### 6. 接口（Interface）
```go
// 定义接口
type Greeter interface {
    Greet() string
}

// 任何实现 Greet() 的类型都隐式满足 Greeter
func Say(g Greeter) {
    fmt.Println(g.Greet())
}

Say(user) // 无需显式声明 implements
```
> 🔑 **鸭子类型**：Go 接口是隐式的，强调“能做什么”而非“是什么”

---

### 7. 并发：Goroutine 与 Channel
```go
// 启动 Goroutine（轻量级线程）
go func() {
    fmt.Println("Running in background")
}()

// Channel 通信
ch := make(chan string)
go func() {
    ch <- "result" // 发送
}()
msg := <-ch // 接收

// Select 多路复用
select {
case msg := <-ch:
    fmt.Println(msg)
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
}
```

---

## 第二部分：工程化开发（现代 Go 项目规范）

### 1. 项目初始化（Go Modules）
```bash
# 创建新项目
mkdir myapp && cd myapp
go mod init github.com/yourname/myapp

# 添加依赖
go get github.com/gin-gonic/gin@v1.10.0

# 整理依赖
go mod tidy
```
- `go.mod`：依赖声明
- `go.sum`：依赖校验（防篡改）

---

### 2. 项目结构（Standard Go Project Layout）
```
myapp/
├── cmd/               # 主应用入口
│   └── myapp/
│       └── main.go
├── internal/          # 私有业务代码（不可被外部 import）
│   ├── handler/       # HTTP/gRPC 处理器
│   ├── service/       # 业务逻辑
│   └── repository/    # 数据访问
├── pkg/               # 可复用公共库（可被外部 import）
├── api/               # OpenAPI/Swagger 定义
├── configs/           # 配置文件
├── scripts/           # 构建/部署脚本
├── test/              # E2E 测试
├── go.mod
└── Makefile           # 构建命令
```

> 💡 **关键原则**：
> - `internal/` 下的代码只能被本项目使用
> - 业务逻辑与框架解耦（避免在 `service` 中直接调用 Gin）

---

### 3. 依赖注入（DI）与分层架构
```go
// internal/repository/user.go
type UserRepository struct{ db *sql.DB }
func (r *UserRepository) GetByID(id int) (*User, error) { ... }

// internal/service/user.go
type UserService struct{ repo *UserRepository }
func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.GetByID(id)
}

// cmd/myapp/main.go
func main() {
    db := connectDB()
    repo := &UserRepository{db: db}
    service := &UserService{repo: repo}
    handler := NewUserHandler(service) // 注入 service
    
    r := gin.Default()
    r.GET("/users/:id", handler.GetUser)
    r.Run(":8080")
}
```
> ✅ **优势**：单元测试时可 mock `UserRepository`

---

### 4. 配置管理（Viper + 环境变量）
```go
// configs/config.go
type Config struct {
    Port    int    `mapstructure:"PORT"`
    DBHost  string `mapstructure:"DB_HOST"`
}

func LoadConfig() (*Config, error) {
    v := viper.New()
    v.SetConfigFile(".env")
    v.AutomaticEnv() // 优先读环境变量
    
    var c Config
    if err := v.Unmarshal(&c); err != nil {
        return nil, err
    }
    return &c, nil
}
```
`.env` 文件：
```env
PORT=8080
DB_HOST=localhost
```

---

### 5. 日志（Zap 结构化日志）
```go
// logger/logger.go
func NewLogger() *zap.Logger {
    return zap.New(zapcore.NewCore(
        zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
        zapcore.AddSync(os.Stdout),
        zap.InfoLevel,
    ))
}

// 使用
logger.Info("user created", 
    zap.Int("user_id", user.ID),
    zap.String("name", user.Name))
```
输出：
```json
{"level":"info","msg":"user created","user_id":1,"name":"Alice"}
```

---

### 6. 单元测试（Table-Driven Tests）
```go
// service/user_test.go
func TestUserService_GetUser(t *testing.T) {
    tests := []struct{
        name    string
        id      int
        wantErr bool
    }{
        {"valid user", 1, false},
        {"invalid id", -1, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // mock repo
            mockRepo := &MockUserRepository{}
            service := UserService{repo: mockRepo}
            
            _, err := service.GetUser(tt.id)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```
运行测试：
```bash
go test ./... -v -cover  # 覆盖率报告
```

---

### 7. API 文档（Swagger）
```go
// handler/user.go
// @Summary Get user by ID
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) { ... }
```
生成文档：
```bash
swag init -g cmd/myapp/main.go
```
访问 `http://localhost:8080/swagger/index.html`

---

### 8. 构建与部署（Makefile + Docker）
`Makefile`：
```makefile
.PHONY: build test docker

build:
	CGO_ENABLED=0 go build -o bin/myapp ./cmd/myapp

test:
	go test -race ./...

docker:
	docker build -t myapp .

run: build
	./bin/myapp
```

`Dockerfile`：
```dockerfile
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY bin/myapp .
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["./myapp"]
```

构建镜像：
```bash
make docker
docker run -p 8080:8080 myapp
```

---

## 第三部分：最佳实践清单

| 类别 | 实践 |
|------|------|
| **代码风格** | 遵循 `gofmt` + `golint`，使用 `go vet` 检查 |
| **错误处理** | 包装错误（`fmt.Errorf("...: %w", err)`），避免忽略 err |
| **并发安全** | 使用 `sync.Mutex` 或 `atomic`，避免全局变量 |
| **性能** | 预分配 slice 容量（`make([]T, 0, size)`），避免内存拷贝 |
| **安全** | 不拼接 SQL（用参数化查询），验证输入 |
| **可观测性** | 集成 Prometheus 指标 + Jaeger 链路追踪 |

---

## 学习资源推荐

1. **官方**：
    - [A Tour of Go](https://go.dev/tour/)（交互式教程）
    - [Effective Go](https://go.dev/doc/effective_go)
2. **书籍**：
    - 《Go 语言高级编程》（开源免费）
    - 《Concurrency in Go》
3. **工具链**：
    - [golangci-lint](https://golangci-lint.run/)（代码检查）
    - [Air](https://github.com/cosmtrek/air)（热重载开发）

---

> 💡 **记住 Go 的核心哲学**：  
> **“Less is exponentially more.” — Rob Pike**  
> 简洁 > 复杂，清晰 > 聪明，组合 > 继承

现在，创建你的第一个工程化 Go 项目吧！🚀