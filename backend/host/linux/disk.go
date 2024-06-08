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

type DiskMetric struct {
	WriteKBs  float64 `json:"writeKBs"`
	ReadKBs   float64 `json:"readKBs"`
	WriteIOs  float64 `json:"writeIOs"`
	ReadIOs   float64 `json:"readIOs"`
	WriteTime float64 `json:"writeTime"`
	ReadTime  float64 `json:"readTime"`
}

func (dm *DiskMetric) ToPoint() common.Point {
	return common.Point{
		Measurement: "disk",
		Tags: map[string]string{
			"hostId": os.Getenv("HOST_ID"),
		},
		Fields: map[string]interface{}{
			"writeKBs":  dm.WriteKBs,
			"readKBs":   dm.ReadKBs,
			"writeIOs":  dm.WriteIOs,
			"readIOs":   dm.ReadIOs,
			"writeTime": dm.WriteTime,
			"readTime":  dm.ReadTime,
		},
		Ts: time.Now(),
	}
}

func getDiskStats(device string) ([]int64, error) {
	session, err := common.GetSSHClient().NewSession()
	if err != nil {

		return nil, fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("grep '" + device + " ' /proc/diskstats"); err != nil {
		return nil, fmt.Errorf("failed to run command: %s", err)
	}
	output := b.String()
	fields := strings.Fields(output)
	if len(fields) < 14 {
		return nil, fmt.Errorf("unexpected output: %s", output)
	}

	readIOs, _ := strconv.ParseInt(fields[3], 10, 64)
	readSectors, _ := strconv.ParseInt(fields[5], 10, 64)
	readTime, _ := strconv.ParseInt(fields[6], 10, 64)
	writeIOs, _ := strconv.ParseInt(fields[7], 10, 64)
	writeSectors, _ := strconv.ParseInt(fields[9], 10, 64)
	writeTime, _ := strconv.ParseInt(fields[10], 10, 64)
	return []int64{readIOs, readSectors, readTime, writeIOs, writeSectors, writeTime}, nil
}

func CollectDisk(ctx context.Context, c *common.Collector) {
	device := "vda" // Change this to your device name

	// First set of statistics
	initialStats, err := getDiskStats(device)
	if err != nil {
		c.ErrCh <- err
		return
	}

	time.Sleep(1 * time.Second) // Sleep for 1 second to calculate the rate

	// Second set of statistics
	finalStats, err := getDiskStats(device)
	if err != nil {
		c.ErrCh <- err
		return
	}

	// Calculating the differences
	readIOsDiff := finalStats[0] - initialStats[0]
	writeIOsDiff := finalStats[3] - initialStats[3]
	readSectorsDiff := finalStats[1] - initialStats[1]
	writeSectorsDiff := finalStats[4] - initialStats[4]
	readTimeDiff := finalStats[2] - initialStats[2]
	writeTimeDiff := finalStats[5] - initialStats[5]

	// Calculating the metrics
	readKBs := float64(readSectorsDiff * 512 / 1024)
	writeKBs := float64(writeSectorsDiff * 512 / 1024)

	dm := &DiskMetric{
		WriteKBs:  writeKBs,
		ReadKBs:   readKBs,
		WriteIOs:  float64(writeIOsDiff),
		ReadIOs:   float64(readIOsDiff),
		WriteTime: float64(writeTimeDiff),
		ReadTime:  float64(readTimeDiff),
	}

	select {
	case <-ctx.Done():
		return
	case c.PointCh <- dm.ToPoint():
		return
	}
}
