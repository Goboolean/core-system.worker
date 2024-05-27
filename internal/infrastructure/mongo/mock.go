package mongo

import (
	"context"

	"github.com/Goboolean/common/pkg/resolver"
)

type Mock__StockClient struct {
	stockID   string
	timeSlice string

	data []*StockDocument
}

func Mock__NewStockClient(c *resolver.ConfigMap, data []*StockDocument) (*Mock__StockClient, error) {
	return &Mock__StockClient{
		data: data,
	}, nil
}

// SetTarget sets the target for a specific stock and timeslice.
func (c *Mock__StockClient) SetTarget(stockID string, timeslice string) {
	c.stockID = stockID
	c.timeSlice = timeslice
}

// GetCount gets count of targeted document
func (c *Mock__StockClient) GetCount(ctx context.Context) int {
	return len(c.data)
}

// FindLatestIndexBy finds the index of the most recent document created before or on the given date.
func (c *Mock__StockClient) FindLatestIndexBy(ctx context.Context, timestamp int64) (int, error) {

	if c.data[0].Timestamp > timestamp {
		return 0, nil
	}

	for i := 1; i < len(c.data); i++ {
		if c.data[i-1].Timestamp < timestamp && c.data[i].Timestamp >= timestamp {
			return i, nil
		}
	}

	return len(c.data), nil
}

// ForEachDocument iterates over a range of documents starting from the specified index and executes the given action
func (c *Mock__StockClient) ForEachDocument(ctx context.Context, startIndex int, quantity int, action func(doc StockDocument)) error {
	for i := startIndex; i < startIndex+quantity; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			action(*c.data[i])
		}
	}
	return nil
}

func (c *Mock__StockClient) Ping(ctx context.Context) error {
	return nil
}

// TODO: 커넥션 닫는 부분 구현
func (c *Mock__StockClient) Close() {
	panic("Not Implied")
}
