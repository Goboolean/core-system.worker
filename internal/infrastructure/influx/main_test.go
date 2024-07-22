package influx_test

import (
	"context"
	"os"
	"testing"

	"github.com/Goboolean/core-system.worker/test/container"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	url              = ""
	orderEventBucket = os.Getenv("INFLUXDB_ORDER_EVENT_BUCKET")
	annotationBucket = os.Getenv("INFLUXDB_ANNOTATION_BUCKET")
	token            = os.Getenv("INFLUXDB_TOKEN")
	org              = os.Getenv("INFLUXDB_ORG")
)

var influxC *container.InfluxContainer

var rawInfluxDBClient influxdb2.Client

var (
	productID = "stock.aapl.usa"
	taskID    = "2024-07-05-test"
)

func TestMain(m *testing.M) {
	var err error
	influxC, err = container.InitInfluxContainer(
		context.Background(),
		os.Getenv("INFLUXDB_ORDER_EVENT_BUCKET"),
		os.Getenv("INFLUXDB_ANNOTATION_BUCKET"))
	if err != nil {
		panic(err)
	}

	url = influxC.URL

	rawInfluxDBClient = influxdb2.NewClient(url, token)
	m.Run()
	err = influxC.Terminate(context.Background())
	if err != nil {
		panic(err)
	}
}
