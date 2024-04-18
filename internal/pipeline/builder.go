package pipeline

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/adapter"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
)

type Spec struct {
	FetchJobName     string
	ModelExecJobName string
	AdaptJobName     string
	ResAnalyzeJob    string
	TransmitJobName  string
}

func Build(spec Spec, UserParams *job.UserParams) (*Pipeline, error) {
	fetch, err := fetcher.Create(spec.FetchJobName, UserParams)
	if err != nil {
		return nil, err
	}

	modelExec, err := executer.Create(spec.ModelExecJobName, UserParams)
	if err != nil {
		return nil, err
	}

	adapt, err := adapter.Create(spec.AdaptJobName, UserParams)
	if err != nil {
		return nil, err
	}

	analyze, err := analyzer.Create(spec.ResAnalyzeJob, UserParams)
	if err != nil {
		return nil, err
	}

	transmit, err := transmitter.Create(spec.TransmitJobName, UserParams)
	if err != nil {
		return nil, err
	}

	return newPipeline(
		fetch,
		modelExec,
		adapt,
		analyze,
		transmit,
	)
}
