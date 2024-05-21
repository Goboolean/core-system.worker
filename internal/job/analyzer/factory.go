package analyzer

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

var providerRepo = map[Spec]jobProvider{
	{ID: "stub", InputType: "candlestick"}: func(p *job.UserParams) (Analyzer, error) {
		return NewStub(p)
	},
}

type jobProvider func(p *job.UserParams) (Analyzer, error)

func Create(spec Spec, p *job.UserParams) (Analyzer, error) {

	var provider, ok = providerRepo[spec]
	if !ok {
		return nil, job.ErrNotFoundJob
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create analyze job: %w", err)
	}

	return f, nil
}
