//go:build wireinject
// +build wireinject

package fetcher

import (
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/google/wire"
)

type mongoConfig resolver.ConfigMap

func provideMongoConfig() mongoConfig {
	//임시
	return mongoConfig(resolver.ConfigMap{})
}

func provideMongo(c mongoConfig) (*infrastructure.MongoClientStockImpl, error) {
	in := resolver.ConfigMap(c)
	return infrastructure.NewMongoClientStock(&in)
}

func initalizeRealtimeStock(p *job.UserParams) (Fetcher, error) {
	wire.Build(
		provideMongoConfig,
		provideMongo,
		NewRealtimeStock,
		wire.Bind(new(Fetcher), new(*RealtimeStock)),
		wire.Bind(new(infrastructure.MongoClientStock), new(*infrastructure.MongoClientStockImpl)),
	)

	return nil, nil
}

func initalizePastStock(p *job.UserParams) (Fetcher, error) {
	wire.Build(
		provideMongoConfig,
		provideMongo,
		NewPastStock,
		wire.Bind(new(Fetcher), new(*PastStock)),
		wire.Bind(new(infrastructure.MongoClientStock), new(*infrastructure.MongoClientStockImpl)),
	)

	return nil, nil
}
