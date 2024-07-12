package joiner

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

// Joiner is an interface that job implementations for the join stage of the pipeline
// Joiner represents a job interface that performs joining operations.
type Joiner interface {
	job.Common

	// SetRefInput sets the reference input data channel for the joiner.
	SetRefInput(job.DataChan)

	// SetModelInput sets the model input data channel for the joiner.
	SetModelInput(job.DataChan)

	// Output returns the output data channel for the joiner.
	Output() job.DataChan
}
