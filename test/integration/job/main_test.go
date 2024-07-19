package job

import (
	"context"
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/test/container"
)

var (
	url         = ""
	token       = configuration.InfluxDBToken
	org         = configuration.InfluxDBOrg
	tradeBucket = configuration.InfluxDBTradeBucket
)

var influxC *container.InfluxContainer

const (
	testStockID   = "stock.aapl.usa"
	testTimeFrame = "1m"
)

func TestMain(m *testing.M) {
	var err error
	influxC, err = container.InitInfluxContainerWithRandomPort(context.Background(), tradeBucket)
	if err != nil {
		panic(err)
	}

	url = influxC.URL

	m.Run()
	err = influxC.Terminate(context.Background())
	if err != nil {
		panic(err)
	}
}
