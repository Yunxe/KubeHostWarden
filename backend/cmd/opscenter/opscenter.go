package main

import (
	"flag"
	"fmt"
	"kubehostwarden/db"
	mysql "kubehostwarden/db"
	"kubehostwarden/opscenter"

	"github.com/joho/godotenv"
)

func main() {
	var envFilePath string
	flag.StringVar(&envFilePath, "env", "", "Path to .env file")
	flag.Parse()

	if envFilePath != "" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			fmt.Printf("Error loading .env file: %v\n", err)
		} else {
			fmt.Printf("Loaded .env file from %s\n", envFilePath)
		}
	}

	if err := mysql.SetupMysql(); err != nil {
		panic("failed to setup mysql")
	}

	err := db.SetupInfluxDB()
	if err != nil {
		panic(err)
	}

	go func() { opscenter.NewServer() }()

	select {}
}
