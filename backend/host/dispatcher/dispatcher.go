package dispatcher

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/host/common"
	"kubehostwarden/utils/log"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func Dispatch(ctx context.Context, c *common.Collector) {
	writeApi := db.GetInfluxClient().Client.WriteAPI(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-c.ReturnError():
			log.Error("producing error", "error", err.Error(), "host", os.Getenv("SSH_HOST"), "type", c.MetricType)
			continue
		case point := <-c.ReturnPonit():
			p := influxdb2.NewPoint(point.Measurement, point.Tags, point.Fields, point.Ts)
			writeApi.WritePoint(p)
			fmt.Printf("write point: %v\n", point)
		}
	}
}
