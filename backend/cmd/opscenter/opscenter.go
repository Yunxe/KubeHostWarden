package main

import (
	"flag"
	"kubehostwarden/db"
	"kubehostwarden/opscenter"
	"kubehostwarden/utils/log"

	"github.com/joho/godotenv"
)

func main() {
	defer log.Sync()

	// 定义命令行参数
	var envFilePath string
	flag.StringVar(&envFilePath, "env", "", "Path to .env file")

	// 解析命令行参数
	flag.Parse()

	// 加载环境变量文件
	if envFilePath != "" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			log.Error("Failed to load .env file", "error", err.Error())
		} else {
			log.Info("Loaded .env file", "path", envFilePath)
		}
	}

	// 设置MySQL数据库
	if err := db.SetupMysql(); err != nil {
		log.Fatal("Failed to setup mysql", "error", err)
		panic(err)
	}

	// 设置InfluxDB
	if err := db.SetupInfluxDB(); err != nil {
		log.Fatal("Failed to setup influxdb", "error", err)
		panic(err)
	}

	// 启动opscenter服务器
	go opscenter.NewServer()

	// 阻止主goroutine退出
	select {}
}
