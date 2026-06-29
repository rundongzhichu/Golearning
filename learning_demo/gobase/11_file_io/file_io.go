// Package file_io 演示 Go 语言中的文件 I/O 操作。
//
// Go 文件 I/O 的特点：
//   - os 包：文件打开、创建、关闭、删除等
//   - io 包：Reader/Writer 接口，提供抽象的 I/O 操作
//   - bufio 包：缓冲 I/O
//   - io/ioutil（已弃用，Go 1.16+ 改为 os.ReadFile / os.WriteFile）
//   - os.ReadFile / os.WriteFile：快捷读写整个文件
//   - defer f.Close()：确保文件关闭
//   - 路径：filepath 包处理路径
//   - embed：编译时嵌入静态文件（Go 1.16+）
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ==================== 文件读写 ====================

func demoReadWriteFile() {
	fmt.Println("=== 文件读写 ===")

	filename := filepath.Join(os.TempDir(), "go_demo_test.txt")

	// 写入整个文件（Go 1.16+）
	data := []byte("Hello, Go I/O!\n第二行内容\n")
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("写入失败: %v\n", err)
		return
	}
	fmt.Printf("写入文件: %s\n", filename)

	// 读取整个文件
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("读取失败: %v\n", err)
		return
	}
	fmt.Printf("内容:\n%s", string(content))

	// 清理
	os.Remove(filename)
}

// ==================== 逐行读取 ====================

func demoReadLineByLine() {
	fmt.Println("\n=== 逐行读取 ===")

	filename := filepath.Join(os.TempDir(), "go_lines.txt")

	// 准备测试文件
	os.WriteFile(filename, []byte("第一行\n第二行\n第三行\n"), 0644)

	// 打开文件
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("打开失败: %v\n", err)
		return
	}
	defer f.Close()

	// 使用 bufio.Scanner 逐行读取
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		fmt.Printf("  行 %d: %s\n", lineNum, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("扫描错误: %v\n", err)
	}

	os.Remove(filename)
}

// ==================== 缓冲读写 ====================

func demoBufferIO() {
	fmt.Println("\n=== 缓冲读写 ===")

	filename := filepath.Join(os.TempDir(), "go_buf.txt")

	// 创建文件
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("创建失败: %v\n", err)
		return
	}
	defer f.Close()
	defer os.Remove(filename)

	// 使用 bufio.Writer
	writer := bufio.NewWriter(f)
	for i := 1; i <= 5; i++ {
		fmt.Fprintf(writer, "行 %d\n", i)
	}
	writer.Flush() // 必须调用 Flush 确保写入

	// 重新定位到开头
	f.Seek(0, 0)

	// 使用 bufio.Reader
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("读取错误: %v\n", err)
			break
		}
		fmt.Printf("  读取: %s", line)
	}
}

// ==================== io.Copy ====================

func demoCopy() {
	fmt.Println("\n=== io.Copy ===")

	srcName := filepath.Join(os.TempDir(), "go_src.txt")
	dstName := filepath.Join(os.TempDir(), "go_dst.txt")

	// 准备源文件
	os.WriteFile(srcName, []byte("复制测试内容\n"), 0644)

	// 打开源文件
	src, err := os.Open(srcName)
	if err != nil {
		fmt.Printf("打开源文件失败: %v\n", err)
		return
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(dstName)
	if err != nil {
		fmt.Printf("创建目标文件失败: %v\n", err)
		return
	}
	defer dst.Close()

	// 复制（一次完成，不必手动循环）
	written, err := io.Copy(dst, src)
	if err != nil {
		fmt.Printf("复制失败: %v\n", err)
		return
	}
	fmt.Printf("复制了 %d 字节\n", written)

	// 验证
	content, _ := os.ReadFile(dstName)
	fmt.Printf("目标内容: %s", string(content))

	os.Remove(srcName)
	os.Remove(dstName)
}

// ==================== Seeker（随机访问） ====================

func demoSeeker() {
	fmt.Println("\n=== Seeker（随机访问） ===")

	filename := filepath.Join(os.TempDir(), "go_seek.txt")
	os.WriteFile(filename, []byte("0123456789"), 0644)
	defer os.Remove(filename)

	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	// 从开头偏移 5 个字节
	f.Seek(5, io.SeekStart)
	buf := make([]byte, 3)
	f.Read(buf)
	fmt.Printf("从开头偏移5: %s\n", string(buf))

	// 从末尾往前偏移 3 个字节
	f.Seek(-3, io.SeekEnd)
	f.Read(buf)
	fmt.Printf("从末尾偏移-3: %s\n", string(buf))
}

// ==================== 文件信息 ====================

func demoFileInfo() {
	fmt.Println("\n=== 文件信息 ===")

	filename := filepath.Join(os.TempDir(), "go_info.txt")
	os.WriteFile(filename, []byte("test"), 0644)
	defer os.Remove(filename)

	// 获取文件信息
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Stat 失败: %v\n", err)
		return
	}

	fmt.Printf("名称: %s\n", info.Name())
	fmt.Printf("大小: %d bytes\n", info.Size())
	fmt.Printf("权限: %s\n", info.Mode())
	fmt.Printf("修改时间: %v\n", info.ModTime())
	fmt.Printf("是否目录: %v\n", info.IsDir())
}

// ==================== 目录操作 ====================

func demoDirectoryOps() {
	fmt.Println("\n=== 目录操作 ===")

	dir := filepath.Join(os.TempDir(), "go_demo_dir")

	// 创建目录
	os.MkdirAll(dir+"/subdir", 0755)
	defer os.RemoveAll(dir)

	// 创建文件
	os.WriteFile(dir+"/file1.txt", []byte("content1"), 0644)
	os.WriteFile(dir+"/file2.txt", []byte("content2"), 0644)
	os.WriteFile(dir+"/subdir/file3.txt", []byte("content3"), 0644)

	// 读取目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("读取目录失败: %v\n", err)
		return
	}
	fmt.Printf("%s 中的条目:\n", dir)
	for _, entry := range entries {
		fmt.Printf("  %s (IsDir=%v)\n", entry.Name(), entry.IsDir())
	}

	// 遍历目录树
	fmt.Println("\n遍历目录树:")
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(dir, path)
		fmt.Printf("  %s (size=%d)\n", relPath, info.Size())
		return nil
	})
}

func main() {
	demoReadWriteFile()
	demoReadLineByLine()
	demoBufferIO()
	demoCopy()
	demoSeeker()
	demoFileInfo()
	demoDirectoryOps()
}
