package common

import (
	"fmt"
	"kubehostwarden/types"
	"kubehostwarden/utils/log"
	"kubehostwarden/utils/sshclient"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var mutex sync.Mutex

type SSH struct {
	Client *ssh.Client
}

var Ssh SSH

func NewSSHConnection() *ssh.Client {
	port, _ := strconv.Atoi(os.Getenv("SSH_PORT"))
	sshInfo := types.SSHInfo{
		EndPoint: os.Getenv("SSH_HOST"),
		Port:     port,
		User:     os.Getenv("SSH_USER"),
		Password: os.Getenv("SSH_PASSWORD"),
	}
	client, err := sshclient.NewSSHClient(sshInfo)
	if err != nil {
		log.Error("Failed to create SSH client", "err:", err)
		return nil
	}
	return client
}

func InitSSHClient() {
	if newConn := NewSSHConnection(); newConn != nil {
		Ssh.Client = newConn
	} else {
		panic("failed to init ssh client")
	}
}

func GetSSHClient() *ssh.Client {
	return Ssh.Client
}

func (sc *SSH) Ping() error {
	session, err := sc.Client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close() // Make sure to close the session
	return nil
}

func ReConnect(maxRetry int) error {
	log.Info("Reconnecting to SSH server","retry times",maxRetry)
	for retry := 0; retry < maxRetry; retry++ {
		client := NewSSHConnection()
		if client == nil {
			fmt.Println("still not working, retry:", retry+1)
			time.Sleep(3 * time.Second)
		} else {
			mutex.Lock()
			if Ssh.Client != nil {
				Ssh.Client.Close()
			}
			Ssh.Client = client
			mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("failed to reconnect to SSH server after %d retries", maxRetry)
}

func HeartBeatDetect() {
	for {
		if err := Ssh.Ping(); err != nil {
			fmt.Println("Failed to ping SSH server", "err", err)
			if err := ReConnect(50); err != nil {
				panic(err)
			} else {
				fmt.Println("ssh server reconnected!")
			}
		}
		time.Sleep(5 * time.Second)
	}
}
