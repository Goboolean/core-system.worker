package jobfactory

import (
	"fmt"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
)

type fetcherProvider func(p *job.UserParams) (fetcher.Fetcher, error)

// wire의 한계 때문에 임시로 이 부분에서 수동 DI
var fetcherProviderRepo = map[fetcher.Spec]fetcherProvider{
	fetcher.Spec{Task: "backtest", ProductType: "stock"}: func(p *job.UserParams) (fetcher.Fetcher, error) {
		//더미
		var mongoCnfig = &resolver.ConfigMap{}
		mongoClient, err := infrastructure.NewMongoClientStock(mongoCnfig)
		if err != nil {
			return nil, err
		}
		return fetcher.NewPastStock(mongoClient, p)
	},
	fetcher.Spec{Task: "realtimeTrade", ProductType: "stock"}: func(p *job.UserParams) (fetcher.Fetcher, error) {
		//더미
		var mongoCnfig = &resolver.ConfigMap{}
		mongoClient, err := infrastructure.NewMongoClientStock(mongoCnfig)
		if err != nil {
			return nil, err
		}
		return fetcher.NewRealtimeStock(mongoClient, p)
	},
}

func CreateFetcher(spec fetcher.Spec, p *job.UserParams) (fetcher.Fetcher, error) {

	var provider, ok = fetcherProviderRepo[spec]
	if !ok {
		return nil, NotFoundJob
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create fetch job: %w", err)
	}

	return f, nil
}
