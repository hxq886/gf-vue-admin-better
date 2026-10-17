// Package main bcrypt 哈希生成工具 (开发辅助, make gen-hash 使用)
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run ./internal/genhash <明文密码>")
		os.Exit(1)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(os.Args[1]), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("生成失败:", err)
		os.Exit(1)
	}
	fmt.Println(string(hash))
}
