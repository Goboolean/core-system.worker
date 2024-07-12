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

// StockTradeCursor is a cursor structure designed for sequentially accessing stock trade data.
// It retrieves an appropriate amount of data from a data source containing past trading data,
// stores it in a buffer, and sequentially provides this data.
type StockTradeCursor struct {
	pastTradeDataSource *influx.DB
	current             time.Time
	limit               int

	// Unique identifier of the product in the format {type}.{ticker}.{locale}
	productID string
	// The interval at which Trade Data is stored.
	// It is also the duration that one Trade Data contains.
	// Follows the following regular expression: ^[0-9]{1,2}[m|s]$.
	// Examples: 1m, 30s, 1s, etc.
	timeFrame           string
	timeFrameDuration   time.Duration
	shouldNotFetchTrade bool

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
		shouldNotFetchTrade: false,
	}, nil
}

// ConfigureStockTradeCursor selects the data that StockTradeCursor will retrieve.
// This function should be called before invoking the Next() function.
func (c *StockTradeCursor) ConfigureStockTradeCursor(startTime time.Time, productID string, timeFrame string) error {
	var err error
	var duration time.Duration
	if duration, err = time.ParseDuration(timeFrame); err != nil {
		return fmt.Errorf("select product: invalid duration format")
	}

	c.current = startTime
	c.productID = productID
	c.timeFrame = timeFrame
	c.timeFrameDuration = duration
	return nil
}

// Next fetches the current trade data pointed by the cursor and moves the cursor to the next trade trade data.
// If an error occurs during data retrieval, it returns (nil, err).
// If there is no more data to retrieve, it returns (nil, nil).
func (c *StockTradeCursor) Next(ctx context.Context) (*model.StockAggregate, error) {

	if len(c.buf) == c.idx && c.shouldNotFetchTrade {
		return nil, nil
	}

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

		if len(buf) < c.limit+1 {
			c.shouldNotFetchTrade = true
		}

		if len(buf) == 0 {
			c.idx = 0
			return nil, nil
		}

		sz := len(buf)
		if len(buf) == c.limit+1 {
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
