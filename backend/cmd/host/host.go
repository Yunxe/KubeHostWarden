package main

import (
	"context"
	"flag"
	"kubehostwarden/db"
	"kubehostwarden/host"
	"kubehostwarden/host/common"
	"kubehostwarden/utils/log"
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
			log.Error("Failed to load .env file", "error", err.Error())
		} else {
			log.Info("Loaded .env file", "path", envFilePath)
		}
	}

	err := db.SetupInfluxDB()
	if err != nil {
		panic(err)
	}

	common.InitSSHClient()
	go common.HeartBeatDetect()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go host.NewHostService(ctx)

	go func() {
		<-signals
		log.Info("Shutting down gracefully...")
		cancel()
	}()

	<-ctx.Done()

	log.Info("Exiting...")
	db.GetInfluxClient().Client.Close()
}
