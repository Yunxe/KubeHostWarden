package common

import (
	"os"
	"time"
)

type Point struct {
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
	Ts          time.Time
}

func GetOSType() string {
	return os.Getenv("SSH_OSTYPE")
}
