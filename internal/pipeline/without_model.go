package pipeline

import (
	"reflect"

	"github.com/Goboolean/core-system.worker/internal/job/adapter"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	"github.com/Goboolean/core-system.worker/internal/util"
)

type WithoutModel struct {
	fetcher     fetcher.Fetcher
	adapter     adapter.Adapter
	analyzer    analyzer.Analyzer
	transmitter transmitter.Transmitter

	demux   *util.ChannelDeMux[error]
	errChan chan error
}

func newWithoutModelWithAdapter(
	fetcher fetcher.Fetcher,
	adapter adapter.Adapter,
	analyzer analyzer.Analyzer,
	transmitter transmitter.Transmitter) (*WithoutModel, error) {

	instance := WithoutModel{
		fetcher:     fetcher,
		adapter:     adapter,
		analyzer:    analyzer,
		transmitter: transmitter,

		demux:   &util.ChannelDeMux[error]{},
		errChan: make(chan error),
	}

	instance.adapter.SetInput(instance.fetcher.Output())
	instance.analyzer.SetInput(instance.adapter.Output())
	instance.transmitter.SetInput(instance.analyzer.Output())

	instance.demux.AddInput(
		instance.fetcher.Error(),
		instance.adapter.Error(),
		instance.analyzer.Error(),
		instance.transmitter.Error(),
	)
	instance.errChan = instance.demux.Output()

	return &instance, nil
}

func newWithoutModelWithoutAdapter(
	fetch fetcher.Fetcher,
	analyze analyzer.Analyzer,
	transmit transmitter.Transmitter) (*WithoutModel, error) {

	instance := WithoutModel{
		fetcher:     fetch,
		analyzer:    analyze,
		transmitter: transmit,

		demux:   &util.ChannelDeMux[error]{},
		errChan: make(chan error),
	}

	instance.analyzer.SetInput(instance.fetcher.Output())
	instance.transmitter.SetInput(instance.analyzer.Output())

	instance.demux.AddInput(
		instance.fetcher.Error(),
		instance.analyzer.Error(),
		instance.transmitter.Error(),
	)
	instance.errChan = instance.demux.Output()

	return &instance, nil
}

func (wom *WithoutModel) Run() {

	wom.fetcher.Execute()
	if !reflect.ValueOf(wom.adapter).IsNil() {
		wom.adapter.Execute()
	}
	wom.analyzer.Execute()
	wom.transmitter.Execute()

}

func (wom *WithoutModel) Stop() {
	wom.fetcher.Stop()
	<-wom.transmitter.Done()
}

func (wom *WithoutModel) Done() chan struct{} {
	return wom.transmitter.Done()
}

func (wom *WithoutModel) Error() chan error {
	return wom.errChan
}
