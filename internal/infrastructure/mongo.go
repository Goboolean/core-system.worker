package infrastructure

import (
	"context"
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

type MongoClientStock struct {
	stockId   string
	timeSlice string
}

func NewMongoClientStock() (*MongoClientStock, error) {
	return &MongoClientStock{}, nil
}

func (c *MongoClientStock) SetStockId(id string) {
	c.stockId = id
}

func (c *MongoClientStock) SetTimeSlice(timeSlice string) {
	c.timeSlice = timeSlice
}

// TODO: 가장 오래된 document가 0번째일 때 startIndex번째부터 endIndex번째 document를 가져오는 부분 구현
func (c *MongoClientStock) FetchItems(ctx context.Context, startIndex int, endIndex int) ([]StockDocument, error) {
	panic("Not Implied")
}

func (c *MongoClientStock) GetCount() (int, error) {
	panic("Not Implied")
}

func (c *MongoClientStock) ping(ctx context.Context) error {
	panic("Not Implied")
}

// TODO: 커넥션 닫는 부분 구현
func (c *MongoClientStock) Close() {
	panic("Not Implied")
}
