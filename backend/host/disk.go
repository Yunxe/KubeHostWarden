package host

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Disk struct {
	HostId     string  `json:"host_id"`
	Tps        float64 `json:"tps"`
	KBPerTrans float64 `json:"KB_Per_Transaction"`
	MBPerSec   float64 `json:"MB_Per_Second"`
}

func (Disk) Name() string {
	return "disk"
}

func (c *Collector) DoCollectDisk() {
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
				var disk *Disk
				switch c.os {
				case "darwin":
					disk = c.DoCollectDiskDarwin()
				}
				if disk != nil {
					c.diskDataCh <- disk
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
			case disk := <-c.diskDataCh:
				p = influxdb2.NewPoint(
					"disk",
					map[string]string{"host_id": disk.HostId},
					map[string]interface{}{
						"tps":                disk.Tps,
						"KB_Per_Transaction": disk.KBPerTrans,
						"MB_Per_Second":      disk.MBPerSec,
					},
					time.Now(),
				)
				c.writeApi.WritePoint(p)
				fmt.Printf("write disk data: %v\n", disk)
			}
		}
	}()
}

func (c *Collector) DoCollectDiskDarwin() *Disk {
	session, err := c.client.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: ", err)
		return nil
	}
	defer session.Close()

	output, err := session.CombinedOutput("iostat -d -w 1 -c 2 | tail -n 1")
	if err != nil {
		fmt.Println("Failed to execute command: ", err)
		return nil
	}

	rawData := strings.TrimSpace(string(output))
	matches := strings.Fields(rawData)

	tps,_:= strconv.ParseFloat(matches[0], 64)
	KBPerTrans,_:= strconv.ParseFloat(matches[1], 64)
	MBPerSec,_:= strconv.ParseFloat(matches[2], 64)

	disk := &Disk{
		HostId:     c.hostId,
		Tps:        tps,
		KBPerTrans: KBPerTrans,
		MBPerSec:   MBPerSec,
	}
	return disk
}
