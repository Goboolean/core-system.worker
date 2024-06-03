package fetcher

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Fetcher interface {
	job.Common

	Output() job.DataChan
}

// FetchingSessin represents a session for fetching data points.
type FetchingSessin interface {
	// TryNext attempts to move to the next data point in the session.
	// It returns true if there is a next data point, false otherwise.
	TryNext() bool

	// Next returns the next FetchingSessin and an error, if any.
	Next() (FetchingSessin, error)

	// DecodedValue decodes the value of the data point into the provided variable.
	// It returns an error if the decoding fails.
	DecodedValue(v any) error

	Close() error
}

// Fetcher is an interface that defines methods for fetching data from InfluxDB.
type TradeReposityPort interface {

	// SelectProduct selects a product by its ID and time frame.
	SelectProduct(ID string, timeFrame string, productType string)

	// SetRangeByTime sets the time range for fetching data.
	SetRangeByTime(fromTimestamp int64, toTimestamp int64)

	// Session returns a fetching session for executing the fetch operation.
	Session() FetchingSessin
}

// Fetcher is an interface that defines methods for fetching data from InfluxDB.
type TradeStreamPort interface {

	// SelectProduct selects a product by its ID and time frame.
	SelectProduct(ID string, timeFrame string, productType string)

	// Session returns a fetching session for executing the fetch operation.
	Session() FetchingSessin
}
