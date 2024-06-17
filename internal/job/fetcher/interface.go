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

// FetchingSession represents a session to fetch trade data in order
// FetchingSession is used to iterate over
// selected range of trade repository or trade stream
type FetchingSession interface {
	// Next advances the session to the next item.
	// It returns true if there is a next item, false otherwise.
	//
	// Next MUST be called initially to retrieve the first item.
	Next() bool

	// Value returns the current item in the fetching session.
	// It returns the item and an error if there was an error retrieving the item.
	Value(ctx context.Context) (any, error)
}

type TradeRepository interface {
	// SelectProduct selects a product by ID, time frame, and product type.
	SelectProduct(ID string, timeFrame string, productType string)

	// SetRangeByTime sets the time range for trade data.
	SetRangeByTime(from time.Time, to time.Time)

	// SetRangeAll sets the range to retrieve all available data entries without any specific limit.
	SetRangeAll()

	// SetRangeByNumberAndEndTime sets the range to retrieve a specified number of trade data
	// that were created just before the given end time,
	SetRangeByNumberAndEndTime(num int, end time.Time)

	// Session returns a fetching session.
	// Before you call session, you must select product and set range
	Session() (FetchingSession, error)

	// Close closes the connection
	Close() error
}

type TradeStream interface {
	// SelectProduct selects a product by ID, time frame, and product type.
	SelectProduct(ID string, timeFrame string, productType string)

	// Session returns a fetching session.
	// Before you call session, you must select product.
	Session() (FetchingSession, error)

	// Close closes the connection.
	Close() error
}
