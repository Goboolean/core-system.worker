package entity

import (
	"context"
	"io"

	"github.com/Goboolean/worker/internal/domain/port/out"
	"github.com/Goboolean/worker/internal/domain/vo"
	log "github.com/sirupsen/logrus"
)




type Model struct {
	data  <-chan *vo.StockAggregate

	stdin io.WriteCloser
	stdout io.ReadCloser

	session out.ModelSession
}


func NewModel(ctx context.Context, stockId string) (*Model, error) {
	return &Model{}, nil
}


// Multiple provider can be set.
func (m *Model) SetDataProvider(ch <-chan *vo.StockAggregate) {

	go func() {
		for data := range ch {
			var v interface{} = data
			_, err := m.stdin.Write(v.([]byte))
			if err != nil {
				log.Error(err)
			}
		}
	}()
}

func (m *Model) Result() chan vo.Result {
	return nil
}

func (m *Model) Close() error {
	return m.session.Close()
}