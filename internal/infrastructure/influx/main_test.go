package influx_test

import (
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	url              = configuration.InfluxDBURL
	orderEventBucket = configuration.InfluxDBOrderEventBucket
	annotationBucket = configuration.InfluxDBAnnotationBucket
	token            = configuration.InfluxDBToken
	org              = configuration.InfluxDBOrg
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
