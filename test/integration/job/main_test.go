package job

import (
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
)

var (
	url             = configuration.InfluxDBURL
	token           = configuration.InfluxDBToken
	org             = configuration.InfluxDBOrg
	tradeBucketName = configuration.InfluxDBTradeBucket
)

const (
	testStockID   = "stock.aapl.usa"
	testTimeFrame = "1m"
)

func TestMain(m *testing.M) {
	m.Run()
}
