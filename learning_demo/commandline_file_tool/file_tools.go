package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func listFiles(dir string) error {
	fmt.Printf("Listing files in: %s\n", dir)
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 跳过根目录本身（只显示内容）
		if path == dir {
			return nil
		}
		fmt.Println(path)
		return nil
	})
}

func removeFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to remove %s: %v", path, err)
	}
	fmt.Printf("Removed: %s\n", path)
	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  filetool list <directory>    - List all files in the directory")
	fmt.Println("  filetool rm <file>           - Remove a file")
	fmt.Println("  filetool help                - Show this help message")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		if len(os.Args) < 3 {
			fmt.Println("Error: missing directory argument for 'list'")
			printUsage()
			os.Exit(1)
		}
		dir := os.Args[2]
		if err := listFiles(dir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "rm", "remove":
		if len(os.Args) < 3 {
			fmt.Println("Error: missing file argument for 'rm'")
			printUsage()
			os.Exit(1)
		}
		filePath := os.Args[2]
		if err := removeFile(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}
