package fetcher

import (
	"context"
	"fmt"
	"time"

	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/fetch-system.IaC/pkg/influx"
	infraModel "github.com/Goboolean/fetch-system.IaC/pkg/model"
	"github.com/cenkalti/backoff"
)

// StockTradeCursor는 특정 시점 이후 StockTradeData를 1개씩 가져올 수 있도록 하는
type StockTradeCursor struct {
	pastTradeDataSource *influx.DB
	current             time.Time
	limit               int

	productID         string
	timeFrame         string
	timeFrameDuration time.Duration

	buf []*model.StockAggregate
	idx int
}

const DefaultLimit = 100

func NewStockTradeCursor(dataSource *influx.DB) (*StockTradeCursor, error) {
	return &StockTradeCursor{
		pastTradeDataSource: dataSource,
		limit:               DefaultLimit,
		buf:                 make([]*model.StockAggregate, 0),
		idx:                 0,
	}, nil
}

func (c *StockTradeCursor) SetStartTime(t time.Time) {
	c.current = t
}

func (c *StockTradeCursor) SelectProduct(productID string, timeFrame string) error {

	var err error
	var duration time.Duration
	if duration, err = time.ParseDuration(timeFrame); err != nil {
		return fmt.Errorf("select product: invalid duration format")
	}
	c.productID = productID
	c.timeFrame = timeFrame
	c.timeFrameDuration = duration
	return nil
}

func (c *StockTradeCursor) Next(ctx context.Context) (*model.StockAggregate, error) {
	if c.idx >= len(c.buf) {
		b := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
		var buf []*infraModel.StockAggregate
		if err := backoff.Retry(func() error {
			var err error
			buf, err = c.pastTradeDataSource.FetchLimitedTradeAfter(ctx, c.productID, c.timeFrame, c.current, c.limit+1)
			if err != nil {
				return err
			}
			return nil
		}, b); err != nil {
			return nil, err
		}

		sz := 0
		if len(buf) > 0 {
			sz = len(buf) - 1
		}
		c.buf = make([]*model.StockAggregate, sz)
		c.current = buf[len(buf)-1].Time

		//discard last element of buf to prevent duplication
		for i := 0; i < sz; i++ {
			c.buf[i] = mapStockAggregate(buf[i], c.timeFrameDuration)
		}
		c.idx = 0
	}

	if len(c.buf) == 0 {
		return nil, nil
	}

	data := c.buf[c.idx]
	c.idx++
	return data, nil
}

func mapStockAggregate(in *infraModel.StockAggregate, d time.Duration) *model.StockAggregate {
	return &model.StockAggregate{
		OpenTime:   in.Time.Add(-d).Unix(),
		ClosedTime: in.Time.Unix(),
		Open:       float32(in.Open),
		Close:      float32(in.Close),
		High:       float32(in.High),
		Low:        float32(in.Low),
		Volume:     float32(in.Volume),
	}
}

func (c *StockTradeCursor) Close() error {
	return c.pastTradeDataSource.Close()
}
