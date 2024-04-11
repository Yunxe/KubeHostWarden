package main

import (
	"context"
	"flag"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/host"
	"os"
	"os/signal"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() 
	
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	collectors := host.NewHostCollectors(ctx)
	for _, collector := range collectors {
		go collector.DoCollect() 
	}

	go func() {
		<-signals 
		fmt.Println("Shutting down gracefully...")
		cancel() 
	}()

	<-ctx.Done() 
	for _, collector := range collectors {
		collector.Close() 
	}
}
