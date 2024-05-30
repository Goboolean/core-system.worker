//go:build develop

package fetcher

import "github.com/Goboolean/core-system.worker/internal/job"

var providerRepo = map[Spec]jobProvider{
	{Task: "backtest", ProductType: "stock"}:      initalizePastStock,
	{Task: "realtimeTrade", ProductType: "stock"}: initalizeRealtimeStock,
	{Task: "stub", ProductType: "stock"}: func(p *job.UserParams) (Fetcher, error) {
		return NewStockStub(p)
	},
}
