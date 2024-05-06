package mongo

import "context"

type StockClient interface {
	SetTarget(stockId string, timeslice string)
	// GetCount gets count of targeted document
	GetCount(ctx context.Context) int
	// FindLatestIndexBy finds the index of the most recent document created before or on the given date.
	FindLatestIndexBy(ctx context.Context, timestamp int64) (int, error)
	// ForEachDocument iterates over a range of documents starting from the specified index and executes the given action
	ForEachDocument(ctx context.Context, startIndex int, quantity int, action func(doc StockDocument)) error
}
