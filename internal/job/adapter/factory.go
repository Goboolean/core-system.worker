package adapter

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

type jobProvider func(p *job.UserParams) (Adapter, error)

func Create(spec Spec, p *job.UserParams) (Adapter, error) {

	var provider, ok = providerRepo[spec]
	if !ok {
		return nil, fmt.Errorf("create adapt job: %w", job.ErrNotFoundJob)
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create adapt job: %w", err)
	}

	return f, nil
}
