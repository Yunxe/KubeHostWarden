package host

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHInfo struct {
	Host     string
	Port     int
	User     string
	Password string
}

func Register(info SSHInfo) error {
	err := probe(info)
	if err != nil {
		return err
	}

	return nil
}

func probe(info SSHInfo) error {
	portStr := fmt.Sprintf("%d", info.Port)

	addr := net.JoinHostPort(info.Host, portStr)

	config := &ssh.ClientConfig{
		User: info.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(info.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("failed to dial: %s", err)

	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("system_profiler SPHardwareDataType"); err != nil {
		return fmt.Errorf("failed to run: %s", err)
	}
	fmt.Println(b.String())
	return nil
}
