package host

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type CPU struct {
	HostId string  `json:"host_id"`
	Idle   float64 `json:"idle"`
	System float64 `json:"system"`
	User   float64 `json:"user"`

	CreatedAt time.Time `json:"created_at"`
}

func (CPU) Name() string {
	return "cpu"
}

func (c *Collector) DoCollectCPU() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(c.frequency)
		defer ticker.Stop()

		for {
			select {
			case <-c.ctx.Done(): 
				return 
			case <-ticker.C:
				var cpu *CPU
				switch c.os {
				case "darwin":
					cpu = c.DoCollectCPUDarwin()
					// 可以添加更多的case来支持不同的操作系统
				}
				if cpu != nil {
					c.cpuDataCh <- cpu
				}
			}
		}
	}()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		var p *write.Point

		errCh := c.writeApi.Errors()
		go func() {
			for {
				select {
				case err := <-errCh:
					fmt.Printf("write error: %s\n", err.Error())
				case <-c.ctx.Done():
					return
				}
			}
		}()

		for {
			select {
			case <-c.ctx.Done(): 
				return 
			case cpuData := <-c.cpuDataCh:
				p = influxdb2.NewPoint("cpu",
					map[string]string{
						"host_id": cpuData.HostId,
					},
					map[string]interface{}{
						"idle":   cpuData.Idle,
						"system": cpuData.System,
						"user":   cpuData.User,
					},
					cpuData.CreatedAt,
				)
				c.writeApi.WritePoint(p)
				fmt.Println(cpuData)
			}
		}
	}()
}

func (c *Collector) DoCollectCPUDarwin() *CPU {
	session, err := c.client.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: ", err)
		return nil
	}
	defer session.Close()

	// run command
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("top -l 1 | grep 'CPU usage'"); err != nil {
		fmt.Println("Failed to run: ", err)
		return nil
	}
	output := b.String()

	reCPU := regexp.MustCompile(`CPU usage: (.*)% user, (.*)% sys, (.*)% idle`).FindStringSubmatch(output)
	idle, err := strconv.ParseFloat(reCPU[1], 64)
	if err != nil {
		fmt.Println("Failed to convert idle to float64: ", err)
	}
	system, err := strconv.ParseFloat(reCPU[2], 64)
	if err != nil {
		fmt.Println("Failed to convert system to float64: ", err)
	}
	user, err := strconv.ParseFloat(reCPU[3], 64)
	if err != nil {
		fmt.Println("Failed to convert user to float64: ", err)
	}

	if reCPU != nil {
		cpu := &CPU{
			Idle:      idle,
			System:    system,
			User:      user,
			CreatedAt: time.Now(),
		}
		return cpu
	}
	return nil
}
