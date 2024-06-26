package executer

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

// ModelExecutor represents an executor for a specific model job.
type ModelExecutor interface {
	job.Common

	// SetInput sets the input data channel for the executor.
	SetInput(job.DataChan)

	// Output returns the output data channel for the executor.
	Output() job.DataChan

	// Cancel notifies the executor to immediately stop all currently running jobs.
	Cancel()
}
