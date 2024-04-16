package mock

import (
	"context"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
)

type Mock__MongoClientStock struct {
	stockId   string
	timeSlice string

	data []*infrastructure.StockDocument
}

func NewMongoClientStock(c *resolver.ConfigMap, data []*infrastructure.StockDocument) (*Mock__MongoClientStock, error) {
	return &Mock__MongoClientStock{
		data: data,
	}, nil
}

// SetTarget sets the target for a specific stock and timeslice.
func (c *Mock__MongoClientStock) SetTarget(stockID string, timeslice string) {
	c.stockId = stockID
	c.timeSlice = timeslice
}

// GetCount gets count of targeted document
func (c *Mock__MongoClientStock) GetCount(ctx context.Context) int {
	return len(c.data)
}

// FindLatestIndexBy finds the index of the most recent document created before or on the given date.
func (c *Mock__MongoClientStock) FindLatestIndexBy(ctx context.Context, timestamp int64) (int, error) {

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
func (c *Mock__MongoClientStock) ForEachDocument(ctx context.Context, startIndex int, quantity int, action func(doc infrastructure.StockDocument)) error {
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

func (c *Mock__MongoClientStock) Ping(ctx context.Context) error {
	return nil
}

// TODO: 커넥션 닫는 부분 구현
func (c *Mock__MongoClientStock) Close() {
	panic("Not Implied")
}
