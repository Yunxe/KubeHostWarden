package db

import (
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type influxdbConfig struct {
	token  string
	url    string
}

type InfluxClient struct {
	Client influxdb2.Client
}

var ic InfluxClient

func GetInfluxClient() InfluxClient {
	return ic
}

func SetupInfluxDB() error {
	if GetInfluxClient().Client != nil {
		return nil
	}

	config := &influxdbConfig{
		token:  os.Getenv("INFLUXDB_TOKEN"),
		url:    os.Getenv("INFLUXDB_URL"),
	}

	client := influxdb2.NewClient(config.url, config.token)

	ic = InfluxClient{Client: client}

	return nil
}
