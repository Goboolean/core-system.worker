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

type PastStock struct {
	timeSlice           string
	isFetchingFullRange bool
	startTime           time.Time // Unix timestamp of start time
	endTime             time.Time
	stockID             string

	cursor *StockTradeCursor

	out job.DataChan `type:"*StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.

	stop *util.StopNotifier
}

func NewPastStock(stockCursor *StockTradeCursor, parmas *job.UserParams) (*PastStock, error) {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.

	instance := &PastStock{
		timeSlice:           DefaultTimeSlice,
		isFetchingFullRange: DefaultIsFetchingFullRange,
		cursor:              stockCursor,
		stop:                util.NewStopNotifier(),
		out:                 make(job.DataChan),
	}

	if !parmas.IsKeyNilOrEmpty(job.ProductID) {

		val, ok := (*parmas)[job.ProductID]
		if !ok {
			return nil, fmt.Errorf("create past stock fetch job: %w", ErrInvalidStockID)
		}

		instance.stockID = val
	}

	if !parmas.IsKeyNilOrEmpty(job.StartDate) {
		instance.isFetchingFullRange = false

		val, err := strconv.ParseInt((*parmas)[job.StartDate], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("create past stock fetch job: %w", err)
		}

		instance.startTime = time.Unix(val, 0)

	}

	if !parmas.IsKeyNilOrEmpty(job.EndDate) {
		instance.isFetchingFullRange = false

		val, err := strconv.ParseInt((*parmas)[job.EndDate], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("create past stock fetch job: %w", err)
		}

		instance.endTime = time.Unix(val, 0)

	}

	return instance, nil
}

func (ps *PastStock) Execute() error {

	defer close(ps.out)
	defer ps.cursor.Close()

	ctx, cancel := context.WithCancel(context.Background())

	//stop sig를 받았을 때 하던 작업을 멈추고 강제종료 하기 위한 부분.
	//graceful shutdown을 원하면 이 부분이 없어도 됩니다.
	go func() {
		<-ps.stop.Done()
		cancel()
	}()

	if err := ps.cursor.ConfigureStockTradeCursor(ps.startTime, ps.stockID, ps.timeSlice); err != nil {
		return fmt.Errorf("execute fetch job:fail to configure trade cursor %w", err)
	}

	e, err := ps.cursor.Next(ctx)
	for ; e != nil; e, err = ps.cursor.Next(ctx) {
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

	return nil
}

func (ps *PastStock) Output() job.DataChan {
	return ps.out
}

func (ps *PastStock) NotifyStop() {
	ps.stop.NotifyStop()
}
