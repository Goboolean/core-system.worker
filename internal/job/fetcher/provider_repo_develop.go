//go:build develop

package fetcher

import "github.com/Goboolean/core-system.worker/internal/job"

var providerRepo = map[Spec]jobProvider{
	{Task: "backtest", ProductType: "stock"}:      initializePastStock,
	{Task: "realtimeTrade", ProductType: "stock"}: initializeRealtimeStock,
	{Task: "stub", ProductType: "stock"}: func(p *job.UserParams) (Fetcher, error) {
		return NewStockStub(p)
	},
}
