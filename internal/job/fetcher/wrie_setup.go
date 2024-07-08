//go:build wireinject
// +build wireinject

package fetcher

import (
	"os"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/fetch-system.IaC/pkg/influx"
	"github.com/google/wire"
)

func provideInfluxConfig() *influx.Opts {
	return &influx.Opts{
		Url:             os.Getenv("INFLUXDB_URL"),
		Token:           os.Getenv("INFLUXDB_TOKEN"),
		Org:             os.Getenv("INFLUXDB_ORG"),
		TradeBucketName: os.Getenv("INFLUXDB_BUCKET"),
	}
}

func InitializePastStock(p *job.UserParams) (Fetcher, error) {
	wire.Build(
		provideInfluxConfig,
		influx.NewDB,
		NewStockTradeCursor,
		NewPastStock,
		wire.Bind(new(Fetcher), new(*PastStock)))
	return &PastStock{}, nil
}
