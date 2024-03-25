package infrastructure

import (
	"context"

	"github.com/Goboolean/common/pkg/resolver"
)

type StockDocument struct {
	Symbol    string
	Open      float32
	Close     float32
	High      float32
	Low       float32
	Average   float32
	Volume    int64
	Timestamp int64
}

type MongoClientStock interface {
	SetTarget(stockId string, timeslice string)
	// GetCount gets count of targeted document
	GetCount(ctx context.Context) int
	// FindLatestIndexBy finds the index of the most recent document created before or on the given date.
	FindLatestIndexBy(ctx context.Context, timestamp int64) (int, error)
	// ForEachDocument iterates over a range of documents starting from the specified index and executes the given action
	ForEachDocument(ctx context.Context, startIndex int, quantity int, action func(doc StockDocument)) error
}

type MongoClientStockImpl struct {
	stockId   string
	timeSlice string
}

func NewMongoClientStock(c *resolver.ConfigMap) (*MongoClientStockImpl, error) {
	return &MongoClientStockImpl{}, nil
}

// SetTarget sets the target for a specific stock and timeslice.
func (c *MongoClientStockImpl) SetTarget(stockID string, timeslice int) {

}

// GetCount gets count of targeted document
func (c *MongoClientStockImpl) GetCount(ctx context.Context) int {
	panic("Not Implied")
}

// FindLatestIndexBy finds the index of the most recent document created before or on the given date.
func (c *MongoClientStockImpl) FindLatestIndexBy(ctx context.Context, timestamp int64) (int, error) {
	panic("Not Implied")
}

// ForEachDocument iterates over a range of documents starting from the specified index and executes the given action
func (c *MongoClientStockImpl) ForEachDocument(ctx context.Context, startIndex int, quantity int, action func(doc StockDocument)) error {
	panic("Not Implied")
}

func (c *MongoClientStockImpl) Ping(ctx context.Context) error {
	panic("Not Implied")
}

// TODO: 커넥션 닫는 부분 구현
func (c *MongoClientStockImpl) Close() {
	panic("Not Implied")
}
