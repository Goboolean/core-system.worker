package transmitter

import (
	"context"

	"github.com/Goboolean/core-system.worker/internal/job"
)

type Transmitter interface {
	job.Common

	SetInput(job.DataChan)
}

// Sender is an interface for sending data to InfluxDB.
type Sender interface {

	// CreateCollection creates a new collection in InfluxDB with the given name.
	CreateCollection(name string)

	// AsyncWrite asynchronously writes the given event to InfluxDB.
	AsyncWrite(data any)

	// Flush flushes any pending writes to InfluxDB.
	Flush(ctx context.Context)

	// Close closes the connection to InfluxDB.
	Close() error
}
