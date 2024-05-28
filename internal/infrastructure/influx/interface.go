package influx

import (
	"context"
)

type Point map[string]interface{}

// FetchingSessin represents a session for fetching data from InfluxDB.
type FetchingSessin interface {
	// TryNext attempts to move to the next data point in the session.
	// It returns true if there is a next data point, false otherwise.
	TryNext() bool

	// Next returns the next FetchingSessin and an error, if any.
	Next() (FetchingSessin, error)

	// Value returns the value of type Point and an error, if any.
	Value() (Point, error)
}

// Fetcher is an interface that defines methods for fetching data from InfluxDB.
type Fetcher interface {

	// SelectProduct selects a product by its ID and time frame.
	SelectProduct(ID string, timeFrame string)

	// SetRangeByTime sets the time range for fetching data.
	SetRangeByTime(fromTimestamp int64, toTimestamp int64)

	// Session returns a fetching session for executing the fetch operation.
	Session() FetchingSessin
}

// Sender is an interface for sending data to InfluxDB.
type Sender interface {

	// CreateCollection creates a new collection in InfluxDB with the given name.
	CreateCollection(name string)

	// AsyncWrite asynchronously writes the given event to InfluxDB.
	AsyncWrite(event Point)

	// Flush flushes any pending writes to InfluxDB.
	Flush(ctx context.Context)

	// Close closes the connection to InfluxDB.
	Close() error
}
