package main

import (
	"context"
	mysql "kubehostwarden/backend/db"
	"kubehostwarden/backend/opscenter/probe"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	if err:=mysql.SetupMysql();err != nil {
		panic("failed to setup mysql")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "22"
	}
	pint, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	err = probe.Register(context.Background(),probe.SSHInfo{
		Host:     os.Getenv("HOST"),
		Port:     pint,
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		OSType:   os.Getenv("OSTYPE"),
	})
	if err != nil {
		panic(err)
	}

}
