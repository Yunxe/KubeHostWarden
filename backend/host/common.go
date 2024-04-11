package host

import (
	"context"
	"fmt"
	"kubehostwarden/types"
	"net"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"golang.org/x/crypto/ssh"
)

type Collector struct {
	Client   *ssh.Client
	writeApi api.WriteAPI
	OS       string
}

type CollectorAPI interface {
	DoCollectCPUDarwin()
}

func Connect(ctx context.Context, info types.SSHInfo) (*ssh.Client, error) {
	// Connect to ssh endpoint
	portStr := fmt.Sprintf("%d", info.Port)

	addr := net.JoinHostPort(info.EndPoint, portStr)

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
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	return client, nil
}

func NewCollector(client *ssh.Client, os string, writeApi api.WriteAPI) *Collector {
	return &Collector{
		Client:   client,
		OS:       os,
		writeApi: writeApi,
	}
}
