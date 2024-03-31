package main

import (
	"kubehostwarden/backend/opscenter/host"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "22"
	}
	pint, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	err = host.Register(host.SSHInfo{
		Host:     os.Getenv("HOST"),
		Port:     pint,
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

}
