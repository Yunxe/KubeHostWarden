package linux

import (
	"bytes"
	"context"
	"fmt"
	"kubehostwarden/host/common"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CPUMetric struct {
	User   float64 `json:"user"`
	Nice   float64 `json:"nice"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
}

func (cm *CPUMetric) ToPoint() common.Point {
	return common.Point{
		Measurement: "cpu",
		Tags: map[string]string{
			"hostId": os.Getenv("HOST_ID"),
		},
		Fields: map[string]interface{}{
			"idle":   cm.Idle,
			"system": cm.System,
			"user":   cm.User,
			"nice":   cm.Nice,
		},
		Ts: time.Now(),
	}
}

func CollectCPU(ctx context.Context, c *common.Collector) {
	session, err := common.GetSSHClient().NewSession()
	if err != nil {
		c.ErrCh <- fmt.Errorf("failed to create session: %s", err)
		return
	}

	defer session.Close()

	// run command
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("top -b -n 2 | grep 'Cpu(s)' | tail -n 1"); err != nil {
		c.ErrCh <- fmt.Errorf("failed to run command: %s", err)
		return
	}
	output := b.String()

	reCPU := regexp.MustCompile(`%Cpu\(s\):\s+(\d+\.\d+)\s+us,\s+(\d+\.\d+)\s+sy,\s+(\d+\.\d+)\s+ni,\s+(\d+\.\d+)\s+id`).FindStringSubmatch(output)

	user, _ := strconv.ParseFloat(strings.TrimSpace(reCPU[1]), 64)
	system, _ := strconv.ParseFloat(strings.TrimSpace(reCPU[2]), 64)
	nice, _ := strconv.ParseFloat(strings.TrimSpace(reCPU[3]), 64)
	idle, _ := strconv.ParseFloat(strings.TrimSpace(reCPU[4]), 64)

	cm := &CPUMetric{
		Idle:   idle,
		System: system,
		Nice:   nice,
		User:   user,
	}

	select {
	case <-ctx.Done():
		return
	case c.PointCh <- cm.ToPoint():
		return
	}
}
