package dispatcher

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/host/common"
	"kubehostwarden/utils/logger"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// var writeApi api.WriteAPI

// func init() {
// }

func Dispatch(ctx context.Context, c *common.Collector) {
	writeApi := db.GetInfluxClient().Client.WriteAPI(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-c.ReturnError():
			logger.Error("producing error", "error", err.Error(), "host", os.Getenv("HOST"), "type", c.MetricType)
			continue
		case point := <-c.ReturnPonit():
			p := influxdb2.NewPoint(point.Measurement, point.Tags, point.Fields, point.Ts)
			writeApi.WritePoint(p)
			fmt.Printf("write point: %v\n", p)
		}
	}
}