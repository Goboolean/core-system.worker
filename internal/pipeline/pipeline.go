package pipeline

import (
	"context"
	"errors"
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job/adapter"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
)

var ErrTypeNotMatch = errors.New("pipeline: cannot build a pipeline because the types are not compatible between the jobs")

// 아키텍처 설계 상 이 구조는 변경되면 안 된다.
type Pipeline struct {
	fetch      fetcher.Fetcher
	modelExec  executer.ModelExecutor
	adapt      adapter.Adapter
	resAnalyze analyzer.Analyzer
	transmit   transmitter.Transmitter

	ctx    context.Context
	cancel context.CancelFunc
}

func newPipeline(
	fetch fetcher.Fetcher,
	modelExec executer.ModelExecutor,
	adapt adapter.Adapter,
	resAnalyze analyzer.Analyzer,
	transmit transmitter.Transmitter) (*Pipeline, error) {

	return &Pipeline{
		fetch:      fetch,
		modelExec:  modelExec,
		adapt:      adapt,
		resAnalyze: resAnalyze,
		transmit:   transmit,
	}, nil
}

func (p *Pipeline) Run() {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	//TODO: 하나의 컴포넌트에서 발행한 메시지를

	p.modelExec.SetInput(p.fetch.Output())

	p.fetch.Execute()
	p.modelExec.Execute()
	p.adapt.Execute()
	p.resAnalyze.Execute()
	p.transmit.Execute()

}

func (p *Pipeline) Stop() error {

	if err := p.fetch.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown fetch job %w", err)
	}
	if err := p.modelExec.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown model execute job %w", err)
	}
	if err := p.adapt.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown adapt job %w", err)
	}
	if err := p.resAnalyze.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown analyze job %w", err)
	}
	if err := p.transmit.Close(); err != nil {
		return fmt.Errorf("pipeline: failed to shutdown transmit job %w", err)
	}

}
