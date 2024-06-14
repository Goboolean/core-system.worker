package fetcher

//go:generate mockgen -destination=fetching_infra_mock.go -package=fetcher --build_flags=--mod=mod . TradeRepository,FetchingSession,TradeStream

// fetcher는 pipeline의 trade data fetch 단계를 수행할 수 있는 Job을 구현합니다.
// 모든 fetch job은 Fetcher 인터페이스를 구현해야 합니다.

import (
	"context"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
)

// Fetcher is an interface that represents a job fetcher.
type Fetcher interface {
	job.Common

	// Output returns the data channel for the fetched trade data.
	Output() job.DataChan
}

// FetchingSession is an interface that represents a fetching session.
type FetchingSession interface {
	// Next advances the session to the next item.
	// It returns true if there is a next item, false otherwise.
	//
	// Next must be called initially to retrieve the first item.
	Next() bool

	// Value returns the current item in the fetching session.
	// It returns the item and an error if there was an error retrieving the item.
	Value(ctx context.Context) (any, error)
}

// TradeRepository is an interface that represents a trade repository.
type TradeRepository interface {
	// SelectProduct selects a product by ID, time frame, and product type.
	SelectProduct(ID string, timeFrame string, productType string)

	// SetRangeByTime sets the time range for fetching data.
	SetRangeByTime(from time.Time, to time.Time)

	// Session returns a fetching session.
	Session() (FetchingSession, error)

	// Close closes the connection
	Close() error
}

// TradeStream is an interface that represents a trade stream.
type TradeStream interface {
	// SelectProduct selects a product by ID, time frame, and product type.
	SelectProduct(ID string, timeFrame string, productType string)

	// Session returns a fetching session.
	Session() (FetchingSession, error)

	// Close closes the connection.
	Close() error
}
