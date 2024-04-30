package alarm

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/utils/constant"
	"kubehostwarden/utils/email"
	"os"
)

func CheckThreshold(ctx context.Context, to string, hostId string, metric string, subMetric string, threshold float64, thresholdType constant.ThresholdType) {
	queryApi := db.GetInfluxClient().Client.QueryAPI(os.Getenv("INFLUXDB_ORG"))
	query := fmt.Sprintf(`
	from(bucket: "%s")
	|> range(start: -1h)
	|> filter(fn: (r) => r["hostId"] == "%s" and r["_measurement"] == "%s" and r["_field"] == "%s")
	|> last()
	|> yield(name: "last")`, os.Getenv("INFLUXDB_BUCKET"), hostId, metric, subMetric)
	result, err := queryApi.Query(context.Background(), query)
	if err != nil {
		fmt.Println("failed to query data for host: ", err)
		return
	}

	for result.Next() {
		value := result.Record().Value().(float64)
		fmt.Println(value)
		if thresholdType == constant.ABOVE && value >= threshold {
			email.SendEmail(to, "阈值警报", fmt.Sprintf("阈值%s为%s的%s超过%f", metric, hostId, subMetric, threshold))
		} else if thresholdType == constant.BELOW && value < threshold {
			email.SendEmail(to, "阈值警报", fmt.Sprintf("阈值%s为%s的%s低于%f", metric, hostId, subMetric, threshold))
		}

	}

	if result.Err() != nil {
		fmt.Println("query parsing error: ", result.Err())
	}
}
