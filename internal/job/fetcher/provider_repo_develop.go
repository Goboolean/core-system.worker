//go:build develop

package fetcher

import "github.com/Goboolean/core-system.worker/internal/job"

var providerRepo = map[Spec]jobProvider{
	{Task: "backTest", ProductType: "stockStub"}: func(p *job.UserParams) (Fetcher, error) {
		(*p)["numOfGeneration"] = "100"
		return NewStockStub(p)
	},
}
