package reporter

import (
	"context"
	"encoding/json"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/utils/logger"
	"net/http"
	"os"
)

func Report(w http.ResponseWriter, r *http.Request) {
	hostId := r.URL.Query().Get("hostId")
	if hostId == "" {
		http.Error(w, "hostId is required", http.StatusBadRequest)
		return
	}

	data, err := fetchData(hostId)
	if err != nil {
		logger.Error("failed to fetch data", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(data)
	if err != nil {
		logger.Error("failed to marshal data", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

// fetchData 根据 hostId 从 InfluxDB 中查询数据
func fetchData(hostId string) ([]map[string]interface{}, error) {
	queryApi := db.GetInfluxClient().Client.QueryAPI(os.Getenv("INFLUXDB_ORG"))
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -1h)
		|> filter(fn: (r) => r["hostId"] == "%s")
		|> limit(n:30)`, os.Getenv("INFLUXDB_BUCKET"), hostId)

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
