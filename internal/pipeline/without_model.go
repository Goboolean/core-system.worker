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
	fetch    fetcher.Fetcher
	adapt    adapter.Adapter
	analyze  analyzer.Analyzer
	transmit transmitter.Transmitter
}

func newWithoutModelWithAdapter(
	fetch fetcher.Fetcher,
	adapt adapter.Adapter,
	analyze analyzer.Analyzer,
	transmit transmitter.Transmitter) (*WithoutModel, error) {

	instance := WithoutModel{
		fetch:    fetch,
		adapt:    adapt,
		analyze:  analyze,
		transmit: transmit,
	}

	instance.adapt.SetInput(instance.fetch.Output())
	instance.analyze.SetInput(instance.adapt.Output())
	instance.transmit.SetInput(instance.analyze.Output())

	return &instance, nil
}

func ewWithoutModelWithoutAdapter(
	fetch fetcher.Fetcher,
	analyze analyzer.Analyzer,
	transmit transmitter.Transmitter) (*WithoutModel, error) {

	instance := WithoutModel{
		fetch:    fetch,
		analyze:  analyze,
		transmit: transmit,
	}

	instance.analyze.SetInput(instance.fetch.Output())
	instance.transmit.SetInput(instance.analyze.Output())

	return &instance, nil
}

func (wom *WithoutModel) Run() {

	wom.fetch.Execute()
	if !reflect.ValueOf(wom.adapt).IsNil() {
		wom.adapt.Execute()
	}
	wom.analyze.Execute()
	wom.transmit.Execute()

}

func (wom *WithoutModel) Stop() error {

	if err := wom.fetch.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown fetch job %w", err)
	}
	if err := wom.adapt.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown adapt job %w", err)
	}
	if err := wom.analyze.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown analyze job %w", err)
	}
	if err := wom.transmit.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown transmit job %w", err)
	}

	return nil
}
