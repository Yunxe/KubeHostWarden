package probe

import (
	"context"
	"fmt"
	mysql "kubehostwarden/backend/db"
	"net"
	"time"

	"github.com/google/uuid"
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
	Id            string `json:"id" gorm:"column:id;primaryKey"`              // 主机ID
	Hostname      string `json:"hostname" gorm:"column:hostname"`             // 主机名
	OS            string `json:"os" gorm:"column:os"`                         // 操作系统名称，例如 macOS, Ubuntu
	OSVersion     string `json:"os_version" gorm:"column:os_version"`         // 操作系统版本
	Kernel        string `json:"kernel" gorm:"column:kernel"`                 // 内核名称，例如 Darwin, Linux
	KernelVersion string `json:"kernel_version" gorm:"column:kernel_version"` // 内核版本
	Arch          string `json:"arch" gorm:"column:arch"`                     // 架构，例如 x86_64, arm64
	IPAddr        string `json:"ip_addr" gorm:"column:ip_addr"`               // 主机的IP地址
	MemoryTotal   string `json:"memory_total" gorm:"column:memory_total"`     // 内存总量，保持统一的单位
	DiskTotal     string `json:"disk_total" gorm:"column:disk_total"`         // 磁盘总量，保持统一的单位

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Host) TableName() string {
	return "host"
}

type probeHelper struct {
	sshClient *ssh.Client
	host      *Host
}

func Register(ctx context.Context, info SSHInfo) error {
	var per probeHelper

	err := per.probe(ctx, info)
	if err != nil {
		return err
	}

	result := mysql.GetClient().Client.WithContext(ctx).Create(per.host)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ph *probeHelper) probe(ctx context.Context, info SSHInfo) error {
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
	ph.sshClient = client
	ph.host = &Host{}

	switch info.OSType {
	case "darwin":
		err := ph.probeDarwin(ctx)
		if err != nil {
			return err
		}
		ph.host.IPAddr = info.Host
	}
	uuid := uuid.New().String()
	ph.host.Id = uuid
	return nil
}
