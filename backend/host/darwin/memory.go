package darwin

import (
	"bytes"
	"context"
	"fmt"
	"kubehostwarden/host/common"
	"os"
	"regexp"
	"strconv"
	"time"
)

type MemoryMetric struct {
	Used       float64 `json:"used"`
	Wired      float64 `json:"wired"`
	Unused     float64 `json:"unused"`
	Compressed float64 `json:"compressed"`
}

func (mm *MemoryMetric) ToPoint() common.Point {
	return common.Point{
		Measurement: "memory",
		Tags: map[string]string{
			"host": os.Getenv("HOST"),
		},
		Fields: map[string]interface{}{
			"used":       mm.Used,
			"wired":      mm.Wired,
			"unused":     mm.Unused,
			"compressed": mm.Compressed,
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

	// run command
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("top -l 1 -s 0 | head -n 10 | grep 'PhysMem'"); err != nil {
		c.ErrCh <- fmt.Errorf("failed to run command: %s", err)
		return
	}
	output := b.String()

	matches := regexp.MustCompile(`PhysMem: (.*)M used \((.*)M wired, (.*)M compressor\), (.*)M unused`).FindStringSubmatch(output)
	used, _ := strconv.ParseFloat(matches[1], 64)
	wired, _ := strconv.ParseFloat(matches[2], 64)
	unused, _ := strconv.ParseFloat(matches[3], 64)
	compressed, _ := strconv.ParseFloat(matches[4], 64)

	mm := &MemoryMetric{
		Used:       used,
		Wired:      wired,
		Unused:     unused,
		Compressed: compressed,
	}

	select {
	case <-ctx.Done():
		return
	case c.PointCh <- mm.ToPoint():
		return
	}
}
