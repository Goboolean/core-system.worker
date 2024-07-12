package analyzer

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

// jobProvider is a type alias for a function
// that creates a Analyzer based on job.UserParams
type jobProvider func(p *job.UserParams) (Analyzer, error)

// Create generates an appropriate fetcher based on the given spec.
// Create passes userParam to the job during this process
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
