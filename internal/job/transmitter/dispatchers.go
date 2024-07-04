//go:generate mockgen -destination=dispatchers_mock.go -package=transmitter --build_flags=--mod=mod . OrderEventDispatcher,AnnotationDispatcher

package transmitter

import (
	"time"

	"github.com/Goboolean/core-system.worker/internal/model"
)

// OrderEventDispatcher is an interface that represents an order event dispatcher.
type OrderEventDispatcher interface {
	// Dispatch dispatches the given order event.
	Dispatch(taskID string, event *model.OrderEvent)

	// Close closes the dispatcher.
	Close() error
}

// AnnotationDispatcher is an interface that represents an annotation dispatcher.
type AnnotationDispatcher interface {
	// Dispatch dispatches the given data.
	Dispatch(taskID string, data any, createdAt time.Time)

	// Close closes the dispatcher.
	Close() error
}
