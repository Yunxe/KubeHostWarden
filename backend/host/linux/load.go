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

type Load struct {
	One     float64 `json:"one"`
	Five    float64 `json:"five"`
	Fifteen float64 `json:"fifteen"`
}

func (lm *Load) ToPoint() common.Point {
	return common.Point{
		Measurement: "load",
		Tags: map[string]string{
			"hostId": os.Getenv("HOST_ID"),
		},
		Fields: map[string]interface{}{
			"one":     lm.One,
			"five":    lm.Five,
			"fifteen": lm.Fifteen,
		},
		Ts: time.Now(),
	}
}

func CollectLoad(ctx context.Context, c *common.Collector) {
	session, err := common.GetSSHClient().NewSession()
	if err != nil {
		c.ErrCh <- fmt.Errorf("failed to create session: %s", err)
		return
	}

	defer session.Close()

	// run command
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("uptime"); err != nil {
		c.ErrCh <- fmt.Errorf("failed to run command: %s", err)
		return
	}
	output := b.String()
	numberStr := strings.TrimSuffix(strings.Split(output, "load average:")[1], "\n")
	if numberStr == "" {
		c.ErrCh <- fmt.Errorf("failed to get load average")
		return
	}
	matches := strings.Split(numberStr, ",")
	if len(matches) != 3 {
		c.ErrCh <- fmt.Errorf("failed to get load average")
		return
	}

	one, _ := strconv.ParseFloat(matches[0], 64)
	five, _ := strconv.ParseFloat(matches[1], 64)
	fifteen, _ := strconv.ParseFloat(matches[2], 64)

	lm := &Load{
		One:     one,
		Five:    five,
		Fifteen: fifteen,
	}

	select {
	case c.PointCh <- lm.ToPoint():
		return
	case <-ctx.Done():
		return
	}
}
