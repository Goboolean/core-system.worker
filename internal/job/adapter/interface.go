package adapter

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

// Adapter is an interface that Job implementations for the adapt stage of the pipeline
type Adapter interface {
	job.Common

	// SetInput sets the input data channel for the analyzer job.
	SetInput(job.DataChan)

	// Output returns the output data channel for the analyzer job.
	Output() job.DataChan
}

// adapt.Common
// adapt.adaptor
