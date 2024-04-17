package main

import (
	"flag"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/opscenter"

	"github.com/joho/godotenv"
)

func main() {
	// 定义命令行参数
	var envFilePath string
	flag.StringVar(&envFilePath, "env", "", "Path to .env file")

	// 解析命令行参数
	flag.Parse()

	// 加载环境变量文件
	if envFilePath != "" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			fmt.Printf("Error loading .env file: %v\n", err)
		} else {
			fmt.Printf("Loaded .env file from %s\n", envFilePath)
		}
	}

	// 设置MySQL数据库
	if err := db.SetupMysql(); err != nil {
		panic("failed to setup mysql")
	}

	// 设置InfluxDB
	if err := db.SetupInfluxDB(); err != nil {
		panic(err)
	}

	// 启动opscenter服务器
	go opscenter.NewServer()

	// 阻止主goroutine退出
	select {}
}
