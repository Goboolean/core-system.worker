package executer

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

type jobProvider func(p *job.UserParams) (ModelExecutor, error)

var providerRepo = map[Spec]jobProvider{}

func Create(spec Spec, p *job.UserParams) (ModelExecutor, error) {

	var provider, ok = providerRepo[spec]
	if !ok {
		return nil, fmt.Errorf("create model execute job: %w", job.ErrNotFoundJob)
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create model execute job: %w", err)
	}

	return f, nil
}
