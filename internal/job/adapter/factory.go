package adapter

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

type jobProvider func(p *job.UserParams) (Adapter, error)

var providerRepo = map[Spec]jobProvider{
	Spec{InputType: "candlestick", OutputType: "candlestick"}: func(p *job.UserParams) (Adapter, error) {
		return Dummy{}, nil
	},
}

func Create(spec Spec, p *job.UserParams) (Adapter, error) {

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
