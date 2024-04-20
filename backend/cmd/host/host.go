package main

import (
	"context"
	"flag"
	"kubehostwarden/db"
	"kubehostwarden/host"
	"kubehostwarden/utils/logger"
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
			logger.Error("Failed to load .env file", "error", err.Error())
		} else {
			logger.Info("Loaded .env file", "path", envFilePath)
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

	go host.NewHostService(ctx)

	go func() {
		<-signals
		logger.Info("Shutting down gracefully...")
		cancel()
	}()

	<-ctx.Done()

	logger.Info("Exiting...")
	db.GetInfluxClient().Client.Close()
	db.GetInfluxClient().Client.Close()
}
