// retrieve.go
package reporter

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Retrieve(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataCh := make(chan map[string]interface{})
	go DataStream(ctx, dataCh)

	for data := range dataCh {
		if err := conn.WriteJSON(data); err != nil {
			log.Println("Error sending data via websocket:", err)
			break
		}
	}
}

// DataStream 向外部提供一个数据通道
func DataStream(ctx context.Context, dataCh chan<- map[string]interface{}) {
	queryApi := db.GetInfluxClient().Client.QueryAPI(os.Getenv("INFLUXDB_ORG"))
	ticker := time.NewTicker(5 * time.Second) // 每5秒查询一次

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(dataCh)
			return
		case <-ticker.C:
			query := fmt.Sprintf(`from(bucket: "%s")
				|> range(start: -1h)
				|> filter(fn: (r) => r._measurement == "cpu" or r._measurement == "memory" or r._measurement == "disk" or r._measurement == "load")
				|> limit(n:10)`, os.Getenv("INFLUXDB_BUCKET"))

			result, err := queryApi.Query(context.Background(), query)
			if err != nil {
				log.Printf("Error querying data: %s", err)
				continue
			}

			for result.Next() {
				dataCh <- result.Record().Values()
			}

			if result.Err() != nil {
				log.Printf("Query parsing error: %s", result.Err().Error())
			}
		}
	}
}
