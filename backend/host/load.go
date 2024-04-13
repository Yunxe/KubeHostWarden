package host

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Load struct {
	HostId  string  `json:"host_id"`
	One     float64 `json:"one"`
	Five    float64 `json:"five"`
	Fifteen float64 `json:"fifteen"`
}

func (Load) Name() string {
	return "load"
}

func (c *Collector) DoCollectLoad() {
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
				var load *Load
				switch c.os {
				case "darwin":
					load = c.DoCollectLoadDarwin()
				}
				if load != nil {
					c.loadDataCh <- load
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
					fmt.Printf("write error: %s\n", err)
				case <-c.ctx.Done():
					return
				}
			}
		}()

		for {
			select {
			case <-c.ctx.Done():
				return
			case load := <-c.loadDataCh:
				p = influxdb2.NewPoint(
					"load",
					map[string]string{"host_id": c.hostId},
					map[string]interface{}{
						"one":     load.One,
						"five":    load.Five,
						"fifteen": load.Fifteen,
					},
					time.Now(),
				)
				c.writeApi.WritePoint(p)
				fmt.Printf("write load data: %v\n", load)
			}
		}
	}()
}

func (c *Collector) DoCollectLoadDarwin() *Load {
	session, err := c.client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create new session: %s\n", err)
		return nil
	}
	defer session.Close()

	// 获取load数据
	cmd := "uptime | awk -F 'load averages: ' '{print $2}'"
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		fmt.Printf("Failed to get load data: %s\n", err)
		return nil
	}

	// 解析load数据
	matches:= strings.Fields(string(output))

	one, _ := strconv.ParseFloat(matches[0], 64)
	five, _ := strconv.ParseFloat(matches[1], 64)
	fifteen, _ := strconv.ParseFloat(matches[2], 64)

	load := &Load{
		HostId:  c.hostId,
		One:     one,
		Five:    five,
		Fifteen: fifteen,
	}
	return load
}
