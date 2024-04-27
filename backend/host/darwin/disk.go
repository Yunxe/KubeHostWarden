package darwin

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

type DiskMetric struct {
	Tps        float64 `json:"tps"`
	KBPerTrans float64 `json:"KB_Per_Transaction"`
	MBPerSec   float64 `json:"MB_Per_Second"`
}


func (dm *DiskMetric) ToPoint() common.Point {
	return common.Point{
		Measurement: "disk",
		Tags: map[string]string{
			"hostId": os.Getenv("HOST_ID"),
		},
		Fields: map[string]interface{}{
			"tps":        dm.Tps,
			"KBPerTrans": dm.KBPerTrans,
			"MBPerSec":   dm.MBPerSec,
		},
		Ts: time.Now(),
	}
}

func CollectDisk(ctx context.Context, c *common.Collector) {
	session, err := common.GetSSHClient().NewSession()
	if err != nil {
		c.ErrCh <- fmt.Errorf("failed to create session: %s", err)
		return
	}

	defer session.Close()

	// run command
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("iostat -d -w 1 -c 2 | tail -n 1"); err != nil {
		c.ErrCh <- fmt.Errorf("failed to run command: %s", err)
		return
	}
	output := b.String()

	rawData := strings.TrimSpace(string(output))
	matches := strings.Fields(rawData)

	tps, _ := strconv.ParseFloat(matches[0], 64)
	KBPerTrans, _ := strconv.ParseFloat(matches[1], 64)
	MBPerSec, _ := strconv.ParseFloat(matches[2], 64)

	dm := &DiskMetric{
		Tps:        tps,
		KBPerTrans: KBPerTrans,
		MBPerSec:   MBPerSec,
	}

	select {
	case <-ctx.Done():
		return
	case c.PointCh <- dm.ToPoint():
		return
	}
}

