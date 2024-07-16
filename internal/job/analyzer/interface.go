package analyzer

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

// Analyzer is an interface that Job implementations for the analyze stage of the pipeline
type Analyzer interface {
	job.Common

	// SetInput sets the input data channel for the analyzer job.
	SetInput(job.DataChan)

	// Output returns the output data channel for the analyzer job.
	Output() job.DataChan
}
