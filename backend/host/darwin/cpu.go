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

type CPUMetric struct {
	Idle   float64 `json:"idle"`
	System float64 `json:"system"`
	User   float64 `json:"user"`
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
	if err := session.Run("top -l 1 -s 0 | head -n 4 | grep 'CPU usage'"); err != nil {
		c.ErrCh <- fmt.Errorf("failed to run command: %s", err)
		return
	}
	output := b.String()

	reCPU := regexp.MustCompile(`CPU usage: (.*)% user, (.*)% sys, (.*)% idle`).FindStringSubmatch(output)
	idle, _ := strconv.ParseFloat(reCPU[1], 64)
	system, _ := strconv.ParseFloat(reCPU[2], 64)
	user, _ := strconv.ParseFloat(reCPU[3], 64)

	cm := &CPUMetric{
		Idle:   idle,
		System: system,
		User:   user,
	}

	select {
	case <-ctx.Done():
		return
	case c.PointCh <- cm.ToPoint():
		fmt.Printf("cpu point: %v\n", cm.ToPoint())
		return
	}
}
