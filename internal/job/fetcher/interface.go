package fetcher

import (
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
)

// Fetcher is an interface that represents a job fetcher.
type Fetcher interface {
	job.Common

	// Output returns the data channel for the fetched job data.
	Output() job.DataChan
}

// FetchingSession is an interface that represents a fetching session.
type FetchingSession interface {
	// TryNext returns true if there is another session to try.
	TryNext() bool

	// Next returns the next fetching session and an error, if any.
	Next() (FetchingSession, error)

	// DecodedValue decodes the value into the provided variable.
	DecodedValue(v any) error

	// Close closes the fetching session.
	Close() error
}

// TradeRepository is an interface that represents a trade repository.
type TradeRepository interface {
	// SelectProduct selects a product by ID, time frame, and product type.
	SelectProduct(ID string, timeFrame string, productType string)

	// SetRangeByTime sets the time range for fetching data.
	SetRangeByTime(from time.Time, to time.Time)

	// Session returns a fetching session.
	Session() FetchingSession
}

// TradeStream is an interface that represents a trade stream.
type TradeStream interface {
	// SelectProduct selects a product by ID, time frame, and product type.
	SelectProduct(ID string, timeFrame string, productType string)

	// Session returns a fetching session.
	Session() FetchingSession
}
