package mongo

import (
	"context"

	"github.com/Goboolean/common/pkg/resolver"
)

type StockClientImpl struct {
	stockId   string
	timeSlice string
}

func NewStockClientImpl(c *resolver.ConfigMap) (*StockClientImpl, error) {
	return &StockClientImpl{}, nil
}

// SetTarget sets the target for a specific stock and timeslice.
func (c *StockClientImpl) SetTarget(stockID string, timeslice string) {

}

// GetCount gets count of targeted document
func (c *StockClientImpl) GetCount(ctx context.Context) int {
	panic("Not implemented")
}

// FindLatestIndexBy finds the index of the most recent document created before or on the given date.
func (c *StockClientImpl) FindLatestIndexBy(ctx context.Context, timestamp int64) (int, error) {
	panic("Not implemented")
}

// ForEachDocument iterates over a range of documents starting from the specified index and executes the given action
func (c *StockClientImpl) ForEachDocument(ctx context.Context, startIndex int, quantity int, action func(doc StockDocument)) error {
	//별도의 go루틴을 생성하지 말고 동기식으로 구현해주세요
	panic("Not implemented")
}

func (c *StockClientImpl) Ping(ctx context.Context) error {
	panic("Not implemented")
}

// TODO: 커넥션 닫는 부분 구현
func (c *StockClientImpl) Close() {
	panic("Not implemented")
}
