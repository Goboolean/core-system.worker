//go:generate mockgen -destination=dispatcher_mock.go -package=kserve --build_flags=--mod=mod . OrderEventDispatcher AnnotationDispatcher

package transmitter

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

// Transmitter represents a job transmitter that sends data to a specific destination.
type Transmitter interface {
	job.Common

	// SetInput sets the input data channel for the transmitter.
	SetInput(job.DataChan)

	// Done returns a channel that is closed when the transmitter has completed all its tasks
	// and cleaned up the given infrastructure.
	Done() chan struct{}
}
