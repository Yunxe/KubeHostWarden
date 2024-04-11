package probe

import (
	"context"
	"encoding/json"
	"fmt"
	mysql "kubehostwarden/db"
	"kubehostwarden/types"
	"net/http"
	"time"

	"golang.org/x/crypto/ssh"
)

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
	sshInfo   types.SSHInfo
	sshClient *ssh.Client
	host      *Host
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var info types.SSHInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request body"}`))
		return
	}

	var pHelper probeHelper

	pHelper.sshInfo = info
	ctx := context.Background()

	err = pHelper.probe(ctx, info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"failed to register new host,error": "%s"}`, err)))
		return
	}

	result := mysql.GetMysqlClient().Client.WithContext(ctx).Create(pHelper.host)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"failed to insert into db,error": "%s"}`, result.Error)))
		return
	}

	err = pHelper.createPod(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"failed to create pod,error": "%s"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
}
