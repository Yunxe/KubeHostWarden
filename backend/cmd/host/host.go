package main

import (
	"context"
	"flag"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/host"
	"kubehostwarden/types"
	"os"
	"strconv"

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

	err := db.SetupInfluxDB()
	if err != nil {
		panic(err)
	}
	writeApi := db.GetInfluxClient().Client.WriteAPI(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "22"
	}
	pint, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	client, err := host.Connect(context.Background(), types.SSHInfo{
		EndPoint: os.Getenv("HOST"),
		Port:     pint,
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		OSType:   os.Getenv("OSTYPE"),
	})
	if err != nil {
		panic(err)
	}
	collector := host.NewCollector(client, os.Getenv("OSTYPE"), writeApi)

	collector.DoCollectCPU()

	select {}
}
