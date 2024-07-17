package fetcher

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

var (
	ErrInvalidStockID       = errors.New("fetch: can't parse stockID")
	ErrDocumentTypeMismatch = errors.New("fetch: mongo: document type mismatch")
)

// PastStock retrieves past stock trade data sequentially and wraps each piece into a model.Packet,
// then sends it to the output channel.
// PastStock fetches data for the specified stock trade data one at a time within the given range.
type PastStock struct {
	timeFrame string
	startTime time.Time
	endTime   time.Time
	stockID   string

	cursor *StockTradeCursor

	out job.DataChan `type:"*StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.

	stop *util.StopNotifier
}

// Parameter List:
// job.ProductID: The unique identifier of the product in the format {type}.{ticker}.{locale}.
// job.StartDate: The start date for data collection.
// job.EndDate: The end date for data collection.
// job.TimeFrame: The interval at which Trade Data is stored.
func NewPastStock(stockCursor *StockTradeCursor, parmas *job.UserParams) (*PastStock, error) {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.

	instance := &PastStock{
		timeFrame: DefaultTimeSlice,
		cursor:    stockCursor,
		stop:      util.NewStopNotifier(),
		out:       make(job.DataChan),
	}

	if !parmas.IsKeyNilOrEmpty(job.ProductID) {

		val, ok := (*parmas)[job.ProductID]
		if !ok {
			return nil, fmt.Errorf("create past stock fetch job: %w", ErrInvalidStockID)
		}

		instance.stockID = val
	}

	if !parmas.IsKeyNilOrEmpty(job.StartDate) {
		val, err := strconv.ParseInt((*parmas)[job.StartDate], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("create past stock fetch job: %w", err)
		}

		instance.startTime = time.Unix(val, 0)

	}

	if !parmas.IsKeyNilOrEmpty(job.EndDate) {
		val, err := strconv.ParseInt((*parmas)[job.EndDate], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("create past stock fetch job: %w", err)
		}

		instance.endTime = time.Unix(val, 0)

	}

	if !parmas.IsKeyNilOrEmpty(job.TimeFrame) {
		instance.timeFrame = (*parmas)[job.TimeFrame]
	}

	return instance, nil
}

func (ps *PastStock) Execute() error {

	defer close(ps.out)
	defer ps.stop.NotifyStop()
	defer ps.cursor.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// This goroutine stops the ongoing work and forces a shutdown when a termination signal is received.	//graceful shutdown을 원하면 이 부분이 없어도 됩니다.
	go func() {
		<-ps.stop.Done()
		cancel()
	}()

	if err := ps.cursor.ConfigureStockTradeCursor(ps.startTime, ps.stockID, ps.timeFrame); err != nil {
		return fmt.Errorf("execute fetch job:fail to configure trade cursor %w", err)
	}

	for {
		e, err := ps.cursor.Next(ctx)

		select {
		case <-ps.stop.Done():
			return nil
		default:
			if e == nil {
				return nil
			}
			if err != nil {
				return fmt.Errorf("execute fetch job:fail to fetch trade %w", err)
			}

			if e.ClosedTime > ps.endTime.Unix() {
				return nil
			}

			ps.out <- model.Packet{
				Time: time.Unix(e.ClosedTime, 0),
				Data: e,
			}
		}

	}

}

func (ps *PastStock) Output() job.DataChan {
	return ps.out
}

func (ps *PastStock) NotifyStop() {
	ps.stop.NotifyStop()
}
