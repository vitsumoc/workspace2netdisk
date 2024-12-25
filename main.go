package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 需要删除的文件夹名
var dirDelete = []string{"node_modules", ".svn", ".git", "__pycache__", "pip"}

// 需要删除的文件大小 20M
var fileSizeDelete int64 = 20 * 1024 * 1024

// 自动删除不需要的文件夹
var autoDelete = false

func main() {
	// 第一个参数必须是文件夹路径
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalln("请提供一个文件夹路径")
	}

	dirPath := args[0]

	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		log.Fatalf("无法访问路径 %s: %v\n", dirPath, err)
	}

	if !fileInfo.IsDir() {
		log.Fatalf("%s 不是一个文件夹\n", dirPath)
		os.Exit(1)
	}

	// 检查是否有 -autoDelete 标志
	for _, arg := range args {
		if arg == "-autoDelete=true" {
			autoDelete = true
		}
	}

	log.Printf("正在处理 %s", dirPath)

	// 创建 folder.md 记录层级路径 delete.md 记录推荐删除的内容
	folderMd, err := os.Create("folder.md")
	if err != nil {
		log.Fatalf("创建 folder.md 失败: %v\n", err)
	}
	defer folderMd.Close()
	deleteMd, err := os.Create("delete.md")
	if err != nil {
		log.Fatalf("创建 delete.md 失败: %v\n", err)
	}
	defer deleteMd.Close()

	// 用 currentFolderPath 表示当前层文件夹路径
	// 用 currentFolderDeepth 表示当前层文件夹相对于主文件夹深度
	// 递归处理并记录每一层的子文件夹信息
	currentFolderPath := dirPath
	currentFolderDeepth := 0
	readFolder(currentFolderPath, currentFolderDeepth, folderMd, deleteMd)
}

func readFolder(path string, deepth int, fFloder *os.File, fDelete *os.File) {
	// 如果本层是需要删除的文件夹
	if endsWith(path) {
		// 自动删除则不再处理
		if autoDelete {
			err := os.RemoveAll(path)
			if err != nil {
				log.Fatalf("删除文件夹 %s 失败: %v\n", path, err)
			}
			return
		}
		// 否则记录不处理
		fDelete.WriteString(path + "\n")
	}
	// 写入本层记录
	for x := 0; x < deepth; x++ {
		fFloder.WriteString("\t")
	}
	fFloder.WriteString("- ")
	fFloder.WriteString(path + "\n")
	// 处理本层的子项目
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("读取文件夹 %s 失败: %v\n", path, err)
	}

	for _, file := range files {
		subPath := filepath.Join(path, file.Name())
		// 子文件夹则递归
		if file.IsDir() {
			readFolder(subPath, deepth+1, fFloder, fDelete)
		} else {
			// 子文件则记录
			for x := 0; x < deepth; x++ {
				fFloder.WriteString("\t")
			}
			fFloder.WriteString("- ")
			fFloder.WriteString(subPath + "\n")
			// 大文件提示删除
			fInfo, err := file.Info()
			if err != nil {
				log.Fatalf("读取文件 %s 失败: %v\n", subPath, err)
			}
			if fInfo.Size() > fileSizeDelete {
				fDelete.WriteString(subPath + "\n")
			}
		}
	}
}

func endsWith(path string) bool {
	// 检查路径是否以需删除后缀结尾
	for _, end := range dirDelete {
		if strings.HasSuffix(path, end) {
			return true
		}
	}
	return false
}
