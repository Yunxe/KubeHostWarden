package host

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Memory struct {
	HostId     string  `json:"host_id"`
	Used       float64 `json:"used"`
	Wired      float64 `json:"wired"`
	Unused     float64 `json:"unused"`
	Compressed float64 `json:"compressed"`

	CreatedAt time.Time `json:"created_at"`
}

func (Memory) Name() string {
	return "memory"
}

func (c *Collector) DoCollectMemory() {
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
				var memory *Memory
				switch c.os {
				case "darwin":
					memory = c.DoCollectMemoryDarwin()
				}
				if memory != nil {
					c.memoryDataCh <- memory
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
			case memoryData := <-c.memoryDataCh:
				p = influxdb2.NewPoint("memory",
					map[string]string{
						"host_id": memoryData.HostId,
					},
					map[string]interface{}{
						"used":       memoryData.Used,
						"wired":      memoryData.Wired,
						"unused":     memoryData.Unused,
						"compressed": memoryData.Compressed,
					},
					memoryData.CreatedAt,
				)
				c.writeApi.WritePoint(p)
				fmt.Printf("write memory data: %v\n", memoryData)
			}
		}
	}()
}

func (c *Collector) DoCollectMemoryDarwin() *Memory {
	session, err := c.client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %s\n", err)
		return nil
	}
	defer session.Close()

	output, err := session.Output("top -l 1 -s 0 | head -n 8 | grep PhysMem")
	if err != nil {
		fmt.Printf("Failed to get memory info: %s\n", err)
		return nil
	}

	reMemory := regexp.MustCompile(`PhysMem: (.*)M used \((.*)M wired, (.*)M compressor\), (.*)M unused`)
	matches := reMemory.FindStringSubmatch(string(output))

	if len(matches) != 5 { 
		fmt.Println("Failed to parse memory info")
		return nil
	}

	used, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		fmt.Printf("Failed to parse used memory: %s\n", err)
	}
	wired, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		fmt.Printf("Failed to parse wired memory: %s\n", err)
	}
	unused, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		fmt.Printf("Failed to parse unused memory: %s\n", err)
	}
	compressed, err := strconv.ParseFloat(matches[4], 64)
	if err != nil {
		fmt.Printf("Failed to parse compressed memory: %s\n", err)
	}

	return &Memory{
		HostId:     c.hostId,
		Used:       used,
		Wired:      wired,
		Unused:     unused,
		Compressed: compressed,
		CreatedAt:  time.Now(),
	}
}
