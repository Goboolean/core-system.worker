package fetcher

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

type jobProvider func(p *job.UserParams) (Fetcher, error)

func Create(spec Spec, p *job.UserParams) (Fetcher, error) {

	var provider, ok = providerRepo[spec]
	if !ok {
		return nil, fmt.Errorf("create fetch job: %w", job.ErrNotFoundJob)
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create fetch job: %w", err)
	}

	return f, nil
}
