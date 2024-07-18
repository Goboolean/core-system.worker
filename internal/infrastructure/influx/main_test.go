package influx_test

import (
	"os"
	"testing"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	url    = os.Getenv("INFLUXDB_URL")
	bucket = os.Getenv("INFLUXDB_BUCKET")
	token  = os.Getenv("INFLUXDB_TOKEN")
	org    = os.Getenv("INFLUXDB_ORG")
)

var rawInfluxDBClient influxdb2.Client

var (
	productID = "stock.aapl.usa"
	taskID    = "2024-07-05-test"
)

func TestMain(m *testing.M) {
	rawInfluxDBClient = influxdb2.NewClient(url, token)
	m.Run()
	rawInfluxDBClient.Close()
}
