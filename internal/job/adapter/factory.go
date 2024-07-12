package adapter

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

// jobProvider is a type alias for a function
// that creates a Adapter based on job.UserParams
type jobProvider func(p *job.UserParams) (Adapter, error)

// Create generates an appropriate fetcher based on the given spec.
// Create passes userParam to the job during this process
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
