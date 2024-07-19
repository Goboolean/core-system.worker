package influx_test

import (
	"context"
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/test/container"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	url              = ""
	orderEventBucket = configuration.InfluxDBOrderEventBucket
	annotationBucket = configuration.InfluxDBAnnotationBucket
	token            = configuration.InfluxDBToken
	org              = configuration.InfluxDBOrg
)

var influxC *container.InfluxContainer

var rawInfluxDBClient influxdb2.Client

var (
	productID = "stock.aapl.usa"
	taskID    = "2024-07-05-test"
)

func TestMain(m *testing.M) {
	var err error
	influxC, err = container.InitInfluxContainerWithRandomPort(
		context.Background(),
		configuration.InfluxDBOrderEventBucket,
		configuration.InfluxDBAnnotationBucket)
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
