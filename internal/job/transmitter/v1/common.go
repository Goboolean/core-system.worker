package v1

import (
	"errors"
	"fmt"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util/chanutil"
)

type Common struct {
	annotationDispatcher transmitter.AnnotationDispatcher
	orderDispatcher      transmitter.OrderEventDispatcher

	task      model.Task
	productId string

	in job.DataChan
}

var ErrInvalidProductId = errors.New("transmit: can't parse productID")
var ErrInvalidTaskString = errors.New("transmit: can't parse task")

func NewCommon(
	annotationDispatcher transmitter.AnnotationDispatcher,
	orderDispatcher transmitter.OrderEventDispatcher,
	params *job.UserParams,
) (*Common, error) {
	instance := &Common{
		annotationDispatcher: annotationDispatcher,
		orderDispatcher:      orderDispatcher,
	}

	if !params.IsKeyNilOrEmpty(job.ProductID) {

		val, ok := (*params)[job.ProductID]
		if !ok {
			return nil, fmt.Errorf("create past stock fetch job: %w", ErrInvalidProductId)
		}

		instance.productId = val
	}

	if !params.IsKeyNilOrEmpty(job.Task) {

		val, ok := (*params)[job.Task]
		if !ok {
			return nil, fmt.Errorf("create past stock fetch job: %w", ErrInvalidTaskString)
		}

		t, err := model.ParseTask(val)
		if err != nil {
			return nil, fmt.Errorf("create past stock fetch job: %w", err)
		}

		instance.task = t

	}

	return instance, nil

}

func (b *Common) Execute() error {

	defer func() {
		// On successful completion: do nothing
		// On failure: wait until the input channel is closed and consume data in input channel
		// in the goroutine
		go chanutil.DummyChannelConsumer(b.in)
	}()

	defer func() {
		b.annotationDispatcher.Close()
		b.orderDispatcher.Close()
	}()

	//TODO: dispatcher error 처리
	for inPacket := range b.in {
		switch v := inPacket.Data.(type) {
		case *model.TradeCommand:
			b.orderDispatcher.Dispatch(
				&model.OrderEvent{
					ProductID: b.productId,
					Command:   *v,
					CreatedAt: time.Now(),
					Task:      b.task,
				})
		default:
			b.annotationDispatcher.Dispatch(v)
		}
	}

	return nil
}

// SetInput sets the input data channel for the transmitter.
func (b *Common) SetInput(in job.DataChan) {
	b.in = in
}
