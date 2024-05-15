package pipeline

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Goboolean/core-system.worker/internal/job/adapter"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/joinner"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	"github.com/Goboolean/core-system.worker/internal/util"
)

var ErrTypeNotMatch = errors.New("pipeline: cannot build a pipeline because the types are not compatible between the jobs")

// 아키텍처 설계 상 이 구조는 변경되면 안 된다.
type Normal struct {
	//jobs
	fetch      fetcher.Fetcher
	modelExec  executer.ModelExecutor
	adapt      adapter.Adapter
	join       joinner.Joinner
	resAnalyze analyzer.Analyzer
	transmit   transmitter.Transmitter

	//utils
	mux util.ChannelMux

	ctx    context.Context
	cancel context.CancelFunc
}

func newNormalWithAdapter(
	fetch fetcher.Fetcher,
	modelExec executer.ModelExecutor,
	adapt adapter.Adapter,
	join joinner.Joinner,
	resAnalyze analyzer.Analyzer,
	transmit transmitter.Transmitter) (*Normal, error) {

	instance := Normal{
		fetch:      fetch,
		modelExec:  modelExec,
		adapt:      adapt,
		join:       join,
		resAnalyze: resAnalyze,
		transmit:   transmit,

		mux: util.ChannelMux{},
	}

	instance.mux.SetInput(instance.fetch.Output())
	instance.modelExec.SetInput(instance.mux.Output())
	adapt.SetInput(instance.modelExec.Output())
	instance.join.SetModelInput(instance.adapt.Output())
	instance.join.SetRefInput(instance.mux.Output())
	instance.resAnalyze.SetInput(instance.join.Output())
	instance.transmit.SetInput(instance.resAnalyze.Output())

	return &instance, nil
}

func newNormalWithoutAdapter(
	fetch fetcher.Fetcher,
	modelExec executer.ModelExecutor,
	join joinner.Joinner,
	resAnalyze analyzer.Analyzer,
	transmit transmitter.Transmitter) (*Normal, error) {

	instance := Normal{
		fetch:      fetch,
		modelExec:  modelExec,
		join:       join,
		resAnalyze: resAnalyze,
		transmit:   transmit,
	}

	instance.mux.SetInput(instance.fetch.Output())
	instance.modelExec.SetInput(instance.mux.Output())
	instance.join.SetModelInput(instance.modelExec.Output())
	instance.join.SetRefInput(instance.mux.Output())
	instance.resAnalyze.SetInput(instance.join.Output())
	instance.transmit.SetInput(instance.resAnalyze.Output())

	return &instance, nil

}

func (n *Normal) Run() {

	n.fetch.Execute()
	n.modelExec.Execute()
	if !reflect.ValueOf(n.adapt).IsNil() {
		n.adapt.Execute()
	}
	n.join.Execute()
	n.resAnalyze.Execute()
	n.transmit.Execute()
}

func (n *Normal) Stop() error {

	if err := n.fetch.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown fetch job %w", err)
	}
	if err := n.modelExec.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown model execute job %w", err)
	}
	if err := n.adapt.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown adapt job %w", err)
	}
	if err := n.resAnalyze.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown analyze job %w", err)
	}
	if err := n.transmit.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown transmit job %w", err)
	}

	return nil
}
