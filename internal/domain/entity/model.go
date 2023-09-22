package entity

import (
	"context"

	"github.com/Goboolean/core-system.worker/internal/domain/port/out"
	"github.com/Goboolean/core-system.worker/internal/domain/vo"
)

type Model struct {
	data <-chan *vo.StockAggregate

	session out.ModelSession
}

func NewModel(ctx context.Context, stockId string) (*Model, error) {
	return &Model{}, nil
}


func (m *Model) SetDataProvider(sink <-chan *vo.StockAggregate) {

	source := m.session.GetInputChan()

	go func() {
		for data := range sink {
			source <- data
		}
	}()
}


func (m *Model) ResultReceiver() <-chan *vo.Result {
	return m.session.GetOutputChan()
}

func (m *Model) Close() error {
	return m.session.Close()
}
