// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package fetcher

import (

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/fetch-system.IaC/pkg/influx"
)

// Injectors from wire_setup.go:

func InitializePastStock(p *job.UserParams) (Fetcher, error) {
	opts := provideInfluxConfig()
	db, err := influx.NewDB(opts)
	if err != nil {
		return nil, err
	}
	stockTradeCursor, err := NewStockTradeCursor(db)
	if err != nil {
		return nil, err
	}
	pastStock, err := NewPastStock(stockTradeCursor, p)
	if err != nil {
		return nil, err
	}
	return pastStock, nil
}

// wire_setup.go:

func provideInfluxConfig() *influx.Opts {
	return &influx.Opts{
		URL:             configuration.InfluxDBURL,
		Token:           configuration.InfluxDBToken,
		Org:             configuration.InfluxDBOrg,
		TradeBucketName: configuration.InfluxDBTradeBucket,
	}
}
