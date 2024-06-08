package linux

import (
	"bytes"
	"context"
	"fmt"
	"kubehostwarden/host/common"
	"os"
	"strconv"
	"strings"
	"time"
)

type MemoryMetric struct {
	TotalMem     float64 `json:"totalMem"`
	FreeMem      float64 `json:"freeMem"`
	AvailableMem float64 `json:"availableMem"`
	UsedMem      float64 `json:"usedMem"`
	Buffers      float64 `json:"buffers"`
	Cached       float64 `json:"cached"`
	SwapTotal    float64 `json:"swapTotal"`
	SwapFree     float64 `json:"swapFree"`
}

func (mm *MemoryMetric) ToPoint() common.Point {
	return common.Point{
		Measurement: "memory",
		Tags: map[string]string{
			"hostId": os.Getenv("HOST_ID"),
		},
		Fields: map[string]interface{}{
			"totalMem":     mm.TotalMem,
			"freeMem":      mm.FreeMem,
			"availableMem": mm.AvailableMem,
			"usedMem":      mm.UsedMem,
			"buffers":      mm.Buffers,
			"cached":       mm.Cached,
			"swapTotal":    mm.SwapTotal,
			"swapFree":     mm.SwapFree,
		},
		Ts: time.Now(),
	}
}

func CollectMemory(ctx context.Context, c *common.Collector) {
	session, err := common.GetSSHClient().NewSession()
	if err != nil {
		c.ErrCh <- fmt.Errorf("failed to create session: %s", err)
		return
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("cat /proc/meminfo"); err != nil {
		c.ErrCh <- fmt.Errorf("failed to run command: %s", err)
		return
	}
	output := b.String()
	lines := strings.Split(output, "\n")

	memStats := make(map[string]float64)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		value, _ := strconv.ParseFloat(fields[1], 64)
		memStats[key] = value
	}

	totalMem := memStats["MemTotal"]
	freeMem := memStats["MemFree"]
	availableMem := memStats["MemAvailable"]
	buffers := memStats["Buffers"]
	cached := memStats["Cached"]
	swapTotal := memStats["SwapTotal"]
	swapFree := memStats["SwapFree"]
	usedMem := totalMem - freeMem - buffers - cached

	mm := &MemoryMetric{
		TotalMem:     totalMem / 1024, // Convert to MB
		FreeMem:      freeMem / 1024,
		AvailableMem: availableMem / 1024,
		UsedMem:      usedMem / 1024,
		Buffers:      buffers / 1024,
		Cached:       cached / 1024,
		SwapTotal:    swapTotal / 1024,
		SwapFree:     swapFree / 1024,
	}

	select {
	case <-ctx.Done():
		return
	case c.PointCh <- mm.ToPoint():
		return
	}
}
