package fetcher

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

type jobProvider func(p *job.UserParams) (Fetcher, error)

var fetcherProviderRepo = map[Spec]jobProvider{
	{Task: "backtest", ProductType: "stock"}:      initalizePastStock,
	{Task: "realtimeTrade", ProductType: "stock"}: initalizeRealtimeStock,
}

func CreateFetcher(spec Spec, p *job.UserParams) (Fetcher, error) {

	var provider, ok = fetcherProviderRepo[spec]
	if !ok {
		return nil, job.NotFoundJob
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create fetch job: %w", err)
	}

	return f, nil
}
