package types

import "time"

type SSHInfo struct {
	EndPoint string `json:"endpoint" validate:"required"`
	Port     int    `json:"port" validate:"required"`
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
	OSType   string `json:"ostype" validate:"required"`
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

	OwnerId string `json:"owner_id" gorm:"column:owner_id"` // 主机的拥有者ID
	Owner   string `json:"owner" gorm:"column:owner"`       // 主机的拥有者

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Host) TableName() string {
	return "host"
}

type User struct {
	Id       string `json:"id" gorm:"column:id;primaryKey"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Email    string `json:"email" gorm:"column:email"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "user"
}