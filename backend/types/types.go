package types

type SSHInfo struct {
	EndPoint string `json:"endpoint"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	OSType   string `json:"ostype"`
}
