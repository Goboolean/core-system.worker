package job

import (
	"context"
	"os"
	"testing"

	"github.com/Goboolean/core-system.worker/test/container"
)

var (
	url         = ""
	token       = os.Getenv("INFLUXDB_TOKEN")
	org         = os.Getenv("INFLUXDB_ORG")
	tradeBucket = os.Getenv("INFLUXDB_TRADE_BUCKET")
)

var influxC *container.InfluxContainer

const (
	testStockID   = "stock.aapl.usa"
	testTimeFrame = "1m"
)

func TestMain(m *testing.M) {
	var err error
	influxC, err = container.InitInfluxContainer(context.Background(), tradeBucket)
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
