package transmitter

import (
	"context"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

// Transmitter is an interface that represents a job transmitter.
type Transmitter interface {
	job.Common

	// SetInput sets the input data channel for the transmitter.
	SetInput(job.DataChan)
}

// OrderEventDispatcher is an interface that represents an order event dispatcher.
type OrderEventDispatcher interface {
	// Dispatch dispatches the given order event.
	Dispatch(event model.OrderEvent)

	// Flush flushes any pending events in the dispatcher.
	Flush(ctx context.Context)

	// Close closes the dispatcher.
	Close() error
}

// AnnotationDispatcher is an interface that represents an annotation dispatcher.
type AnnotationDispatcher interface {
	// Dispatch dispatches the given data.
	Dispatch(data any)

	// Flush flushes any pending data in the dispatcher.
	Flush(ctx context.Context)

	// Close closes the dispatcher.
	Close() error
}
