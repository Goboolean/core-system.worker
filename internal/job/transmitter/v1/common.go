package v1

import (
	"errors"
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util/chanutil"
)

// Common publishes data to external systems using suitable dispatchers depending on the type of data received from the input channel.
// Common distinguishes incoming data from the input channel into two types: order events and annotations, and dispatches them accordingly.
// Order events are dispatched using the OrderEventDispatcher, while annotations are dispatched using the AnnotationDispatcher.
type Common struct {
	annotationDispatcher transmitter.AnnotationDispatcher
	orderDispatcher      transmitter.OrderEventDispatcher

	task      model.Task
	productId string
	taskID    string

	in job.DataChan
}

var ErrInvalidProductId = errors.New("transmit: can't parse productID")
var ErrInvalidTaskString = errors.New("transmit: can't parse task")

// NewCommon creates new Common instance
//
// Params list
// job.ProductID: The unique identifier of the product in the format {type}.{ticker}.{locale}
// job.Task: The type of task that this application performs
// job.TaskID: The unique identifier of the task that this application performs
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

	if !params.IsKeyNilOrEmpty(job.TaskID) {

		val, ok := (*params)[job.TaskID]
		if !ok {
			return nil, fmt.Errorf("create past stock fetch job: %w", ErrInvalidTaskString)
		}

		instance.taskID = val

	}

	return instance, nil

}

// Execute starts to receive and dispatch data
//
// If the Job fails to perform its task, Execute returns an error.
// If the Job completes successfully, it returns nil.
// DO NOT CALL Execute() TWICE. IT MUST BE PANIC
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
				b.taskID,
				&model.OrderEvent{
					ProductID: b.productId,
					Command:   *v,
					CreatedAt: inPacket.Time,
					Task:      b.task,
				})
		default:
			b.annotationDispatcher.Dispatch(b.taskID, v, inPacket.Time)
		}
	}

	return nil
}

func (b *Common) SetInput(in job.DataChan) {
	b.in = in
}
