//go:generate mockgen -destination=dispatchers_mock.go -package=transmitter --build_flags=--mod=mod . OrderEventDispatcher,AnnotationDispatcher

package transmitter

import (
	"context"

	"github.com/Goboolean/core-system.worker/internal/model"
)

// OrderEventDispatcher is an interface that represents an order event dispatcher.
type OrderEventDispatcher interface {
	// Dispatch dispatches the given order event.
	Dispatch(event *model.OrderEvent)

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
