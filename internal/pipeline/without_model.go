package pipeline

import (
	"reflect"

	"github.com/Goboolean/core-system.worker/internal/job/adapter"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
)

type WithoutModel struct {
	fetcher     fetcher.Fetcher
	adapter     adapter.Adapter
	analyzer    analyzer.Analyzer
	transmitter transmitter.Transmitter
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
	}

	instance.adapter.SetInput(instance.fetcher.Output())
	instance.analyzer.SetInput(instance.adapter.Output())
	instance.transmitter.SetInput(instance.analyzer.Output())

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
	}

	instance.analyzer.SetInput(instance.fetcher.Output())
	instance.transmitter.SetInput(instance.analyzer.Output())

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
}

func (wom *WithoutModel) Done() chan struct{} {
	return wom.transmitter.Done()
}
