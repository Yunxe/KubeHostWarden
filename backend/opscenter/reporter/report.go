package reporter

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	resp "kubehostwarden/utils/responsor"
	"net/http"
	"net/url"
	"os"
)

func Report(ctx context.Context, values url.Values) resp.Responsor {
	hostId := values.Get("hostId")
	metricType := values.Get("mt")
	if hostId == "" || metricType == "" {
		return resp.Responsor{
			Code:    http.StatusBadRequest,
			Message: "hostId and metricType are required",
		}
	}

	data, err := fetchData(hostId, metricType)
	if err != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "data retrieved successfully",
		Result:  data,
	}
}

// fetchData 根据 hostId 从 InfluxDB 中查询数据
func fetchData(hostId string, metricType string) ([]map[string]interface{}, error) {
	queryApi := db.GetInfluxClient().Client.QueryAPI(os.Getenv("INFLUXDB_ORG"))
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -1h)
		|> filter(fn: (r) => r["hostId"] == "%s" and r["_measurement"] == "%s")
		|> limit(n:10)`, os.Getenv("INFLUXDB_BUCKET"), hostId, metricType)

	result, err := queryApi.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query data for host %s: %v", hostId, err)
	}

	var data []map[string]interface{}
	for result.Next() {
		data = append(data, result.Record().Values())
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("query parsing error for host %s: %s", hostId, result.Err().Error())
	}

	return data, nil
}
