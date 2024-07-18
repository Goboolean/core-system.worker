//go:build wireinject
// +build wireinject

package fetcher

import (
	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/fetch-system.IaC/pkg/influx"

	"github.com/google/wire"
)

func provideInfluxConfig() *influx.Opts {
	return &influx.Opts{
		URL:             configuration.InfluxDBURL,
		Token:           configuration.InfluxDBToken,
		Org:             configuration.InfluxDBOrg,
		TradeBucketName: configuration.InfluxDBTradeBucket,
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
