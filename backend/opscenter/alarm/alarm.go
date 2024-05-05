package alarm

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/opscenter/logger"
	"kubehostwarden/utils/constant"
	"kubehostwarden/utils/email"
	"kubehostwarden/utils/log"
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
		log.Error("query error: ", err)
		return
	}

	for result.Next() {
		value := result.Record().Value().(float64)
		fmt.Println(value)
		if thresholdType == constant.ABOVE && value >= threshold {
			logger.Info(ctx.Value(constant.UserIDKey).(string), "阈值警报", "阈值", metric, "主机", hostId, "子指标", subMetric, "超过", threshold)
			email.SendEmail(to, "阈值警报", fmt.Sprintf("阈值%s为%s的%s超过%f", metric, hostId, subMetric, threshold))
		} else if thresholdType == constant.BELOW && value < threshold {
			logger.Info(ctx.Value(constant.UserIDKey).(string), "阈值警报", "阈值", metric, "主机", hostId, "子指标", subMetric, "低于", threshold)
			email.SendEmail(to, "阈值警报", fmt.Sprintf("阈值%s为%s的%s低于%f", metric, hostId, subMetric, threshold))
		}
	}

	if result.Err() != nil {
		log.Error("query error: ", result.Err())
	}
}
