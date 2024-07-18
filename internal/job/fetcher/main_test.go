package fetcher_test

import (
	"os"
	"testing"
)

var (
	url             = os.Getenv("INFLUXDB_URL")
	token           = os.Getenv("INFLUXDB_TOKEN")
	org             = os.Getenv("INFLUXDB_ORG")
	tradeBucketName = os.Getenv("INFLUXDB_TRADE_BUCKET")
)

const (
	testStockID   = "stock.aapl.usa"
	testTimeFrame = "1m"
)

func TestMain(m *testing.M) {
	m.Run()
}
