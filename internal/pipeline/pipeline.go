package pipeline

import (
	"context"
	"errors"

	"github.com/Goboolean/core-system.worker/internal/job/adapt"
	"github.com/Goboolean/core-system.worker/internal/job/analyze"
	"github.com/Goboolean/core-system.worker/internal/job/fetch"
	modelExecute "github.com/Goboolean/core-system.worker/internal/job/model-execute"
	"github.com/Goboolean/core-system.worker/internal/job/transmit"
)

var ErrTypeNotMatch = errors.New("pipeline: cannot build a pipeline because the types are not compatible between the jobs")

// 아키텍처 설계 상 이 구조는 변경되면 안 된다.
type Pipeline struct {
	fetch      fetch.Fetcher
	modelExec  modelExecute.ModelExecutor
	adapt      adapt.Adapter
	resAnalyze analyze.Analyzer
	transmit   transmit.Transmitter

	ctx    context.Context
	cancel context.CancelFunc
}

func newPipeline(fetch, modelExec, adapt, resAnalyze, transmit job.Job) (*Pipeline, error) {

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

	p.fetch.Execute(p.ctx)
	p.modelExec.Execute(p.ctx)
	p.adapt.Execute(p.ctx)
	p.resAnalyze.Execute(p.ctx)
	p.transmit.Execute(p.ctx)

}

func (p *Pipeline) Stop() {
	p.cancel()
}
