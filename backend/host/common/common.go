package common

import (
	"kubehostwarden/types"
	"kubehostwarden/utils/sshclient"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type Point struct {
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
	Ts          time.Time
}

func GetOSType() string {
	return os.Getenv("OSTYPE")
}

var SshClient *ssh.Client
var once sync.Once

func GetSSHClient() *ssh.Client {
	once.Do(func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		sshInfo := types.SSHInfo{
			EndPoint: os.Getenv("HOST"),
			Port:     port,
			User:     os.Getenv("USER"),
			Password: os.Getenv("PASSWORD"),
		}
		client, err := sshclient.NewSSHClient(sshInfo)
		if err != nil {
			log.Fatalf("Failed to create SSH client: %s", err)
		}
		SshClient = client
	})
	return SshClient
}
