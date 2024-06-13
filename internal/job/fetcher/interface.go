package fetcher

// fetcher는 pipeline의 trade data fetch 단계를 수행할 수 있는 Job을 구현합니다.
// 모든 fetch job은 Fetcher 인터페이스를 구현해야 합니다.

import (
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
)

// Fetcher represents a job fetcher that retrieves trade data.
type Fetcher interface {
	job.Common

	// Output returns the data channel for the fetched trade data.
	Output() job.DataChan

	// Stop stops the fetcher and releases any allocated resources.
	Stop()
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
