package analyzer

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

type jobProvider func(p *job.UserParams) (Analyzer, error)

var providerRepo = map[Spec]jobProvider{
	Spec{Id: "abc", InputType: "candlestick"}: func(p *job.UserParams) (Analyzer, error) {
		return Dummy{}, nil
	},
}

func Create(spec Spec, p *job.UserParams) (Analyzer, error) {

	var provider, ok = providerRepo[spec]
	if !ok {
		return nil, job.NotFoundJob
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create analyze job: %w", err)
	}

	return f, nil
}
