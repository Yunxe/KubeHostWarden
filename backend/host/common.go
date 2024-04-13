package host

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	"net"
	"os"
	"sync"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"golang.org/x/crypto/ssh"
)

type Collector struct {
	hostId     string
	metricType string
	os         string //darwin, linux
	
	client     *ssh.Client
	writeApi   api.WriteAPI
	frequency  time.Duration

	cpuDataCh    chan *CPU
	memoryDataCh chan *Memory
	diskDataCh   chan *Disk
	loadDataCh   chan *Load

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewHostCollectors(ctx context.Context) map[string]*Collector {
	var collectorMap = make(map[string]*Collector)

	// 初始化collectorMap
	collectorMap["CPU"] = &Collector{
		cpuDataCh: make(chan *CPU, 5),
	}
	collectorMap["Memory"] = &Collector{
		memoryDataCh: make(chan *Memory, 5),
	}
	collectorMap["Disk"] = &Collector{
		diskDataCh: make(chan *Disk, 5),
	}
	collectorMap["Load"] = &Collector{
		loadDataCh: make(chan *Load, 5),
	}
	// collectorMap["Network"] = &Collector{}

	for mt, collector := range collectorMap {
		// init context
		collector.ctx, collector.cancel = context.WithCancel(ctx)

		// set metric type,host id, os
		collector.metricType = mt
		collector.hostId = os.Getenv("HOST_ID")
		collector.os = os.Getenv("OSTYPE")

		// establish ssh connection
		sshClient, err := NewSSHClient()
		if err != nil {
			fmt.Printf("Failed to init %v ssh client: %s\n", mt, err)
			continue
		}
		collector.client = sshClient

		// establish influxdb write api
		writeApi := db.GetInfluxClient().Client.WriteAPI(
			os.Getenv("INFLUXDB_ORG"),
			os.Getenv("INFLUXDB_BUCKET"),
		)
		collector.writeApi = writeApi

		// set default frequency
		collector.frequency = 5 * time.Second
	}

	return collectorMap
}

func (c *Collector) DoCollect() {
	switch c.metricType {
	case "CPU":
		c.DoCollectCPU()
	case "Memory":
		c.DoCollectMemory()
	case "Disk":
		c.DoCollectDisk()
	case "Load":
		c.DoCollectLoad()
	}
}

func (c *Collector) Close() {
	if c.cancel != nil {
		c.cancel()
	}

	c.wg.Wait()

	if c.client != nil {
		if err := c.client.Close(); err != nil {
			fmt.Printf("Failed to close ssh client: %s\n", err)
		}
	}

	switch c.metricType {
	case "CPU":
		if c.cpuDataCh != nil {
			close(c.cpuDataCh)
		}
	case "Memory":
		if c.memoryDataCh != nil {
			close(c.memoryDataCh)
		}
	}

}

func NewSSHClient() (*ssh.Client, error) {
	// Connect to ssh endpoint
	addr := net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT"))

	config := &ssh.ClientConfig{
		User: os.Getenv("USER"),
		Auth: []ssh.AuthMethod{
			ssh.Password(os.Getenv("PASSWORD")),
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
