package job

import (
	"context"
)

// Job represents a job that can be executed.
type Job interface {
	// Execute executes the job with the given context.
	Execute(ctx context.Context)

	// SetInputChan sets the input channel for the job.
	SetInputChan(chan any)

	// OutputChan returns the output channel for the job.
	OutputChan() chan any

	// returns a channel that will be closed when the job is done.
	Done() chan bool
}
