package pipeline

import (
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

func NewWithoutModelWithAdapter(
	fetcher fetcher.Fetcher,
	adapter adapter.Adapter,
	analyzer analyzer.Analyzer,
	transmitter transmitter.Transmitter) (*WithoutModel, error) {

	instance := WithoutModel{
		fetcher:     fetcher,
		adapter:     adapter,
		analyzer:    analyzer,
		transmitter: transmitter,

		demux:   util.NewChannelDeMux[error](),
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

func NewWithoutModelWithoutAdapter(
	fetch fetcher.Fetcher,
	analyze analyzer.Analyzer,
	transmit transmitter.Transmitter) (*WithoutModel, error) {

	instance := WithoutModel{
		fetcher:     fetch,
		analyzer:    analyze,
		transmitter: transmit,

		demux:   util.NewChannelDeMux[error](),
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

	wom.demux.Execute()
	wom.fetcher.Execute()
	if wom.adapter != nil {
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
