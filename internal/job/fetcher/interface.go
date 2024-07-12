package fetcher

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

// Fetcher is an interface that Job implementations for the fetch stage of the pipeline
type Fetcher interface {
	job.Common

	// Output returns the data channel for the fetched trade data.
	Output() job.DataChan

	// NotifyStop stops the fetcher and releases any allocated resources.
	NotifyStop()
}
