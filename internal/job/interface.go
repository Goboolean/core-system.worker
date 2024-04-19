package job

import (
	"context"
)

// Common represents a job that can be executed.
type Common interface {
	// Execute executes the job with the given context.
	Execute(ctx context.Context)
	Close() error
}
