package manager

import (
	"context"

	"github.com/Goboolean/core-system.worker/internal/domain/entity"
	"github.com/Goboolean/core-system.worker/internal/domain/port/out"
	"github.com/Goboolean/core-system.worker/internal/domain/vo"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	w out.WorkDispatcher

	past out.PastDataFetcher
	real out.RealDataFetcher

	event out.ResultDispatcher
}

func New(ctx context.Context) *Manager {
	return &Manager{}
}

// If run returns error in the middle of running task, it means task is not successfully finished
func (m *Manager) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	// Here are procedure.

	// 1. Subscribe kafka broker's task register event
	// Next, wait for any task to be allocated.

	taskEvent, ok := <-m.w.RegisterWorker()
	if !ok {
		return ErrRegisterFailed
	}

	// 2. Create stock data receiver as a channel.
	// If event type is real, connect to kafka
	// Otherwise if it is past, connect to mongodb

	var channel <-chan *vo.StockAggregate

	switch {
	case taskEvent.Type == vo.Real:
		channel = m.real.GetChannel(ctx, taskEvent.StockId)
		// TODO : connect to kafka
	case taskEvent.Type == vo.Past:
		// TODO : connect to mongodb
		channel = m.past.GetChannel(ctx, taskEvent.ModelId)
	default:
		return ErrInvalidTaskType
	}

	// 3. Create model with it's initializer
	// this is made up of sub procedure.
	// 3-1. Get model as a file from MiniO
	// 3-2. Compile c++ model file to a binary.
	// 3-3. Run a binary and create input channel and output channel

	_model, err := entity.NewModel(ctx, taskEvent.ModelId)
	if err != nil {
		return err
	}

	// 4. Put the data receiver channel to the model input
	_model.SetDataProvider(channel)

	// 5. Put the model output to the kafka message broker

	// 6. Send a message to kafka that model is successfully running
	// Next, wait til the task is finished and free all resources

	for {
		select {
		case <-ctx.Done():
			// Case: model is finished with exit code 1
			if err := context.Cause(ctx); err != nil {
				return err
			}

			// Case: model is finished with exit code 0
			if err := _model.Close(); err == nil {
				return nil
			}

		case result := <-_model.Result():

			ctx := context.WithoutCancel(ctx)
			if err := m.event.SendResult(ctx, result); err != nil {
				log.Error(err)
			}
		}
	}
}
