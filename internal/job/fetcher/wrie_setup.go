//go:build wireinject
// +build wireinject

package fetcher

import (
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/infrastructure/mongo"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/google/wire"
)

type mongoConfig resolver.ConfigMap

func provideMongoConfig() mongoConfig {
	//임시
	return mongoConfig(resolver.ConfigMap{})
}

func provideMongo(c mongoConfig) (*mongo.StockClientImpl, error) {
	in := resolver.ConfigMap(c)
	return mongo.NewStockClientImpl(&in)
}

func initializeRealtimeStock(p *job.UserParams) (Fetcher, error) {
	wire.Build(
		provideMongoConfig,
		provideMongo,
		NewRealtimeStock,
		wire.Bind(new(Fetcher), new(*RealtimeStock)),
		wire.Bind(new(mongo.StockClient), new(*mongo.StockClientImpl)),
	)

	return nil, nil
}

func initializePastStock(p *job.UserParams) (Fetcher, error) {
	wire.Build(
		provideMongoConfig,
		provideMongo,
		NewPastStock,
		wire.Bind(new(Fetcher), new(*PastStock)),
		wire.Bind(new(mongo.StockClient), new(*mongo.StockClientImpl)),
	)

	return nil, nil
}
