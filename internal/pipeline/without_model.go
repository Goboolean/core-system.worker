package pipeline

import (
	"fmt"
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

func ewWithoutModelWithoutAdapter(
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

func (wom *WithoutModel) Stop() error {

	if err := wom.fetcher.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown fetch job %w", err)
	}
	if err := wom.adapter.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown adapt job %w", err)
	}
	if err := wom.analyzer.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown analyze job %w", err)
	}
	if err := wom.transmitter.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown transmit job %w", err)
	}

	return nil
}
