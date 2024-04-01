package host

import (
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
	OSType   string
}

type Host struct {
	Hostname      string   // 主机名
	OS            string   // 操作系统名称，例如 macOS, Ubuntu
	OSVersion     string   // 操作系统版本
	Kernel        string   // 内核名称，例如 Darwin, Linux
	KernelVersion string   // 内核版本
	Arch          string   // 架构，例如 x86_64, arm64
	IPAddr        string   // 主机的IP地址
	MemoryTotal   string   // 内存总量，保持统一的单位
	DiskTotal     string   // 磁盘总量，保持统一的单位
}

type prober struct {
	sshClient *ssh.Client
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
	prober := prober{
		sshClient: client,
	}
	switch info.OSType {
	case "darwin":
		_,err = probeDarwin(prober)
	}
	if err != nil {
		return err
	}
	return nil
}
