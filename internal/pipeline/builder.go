package pipeline

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/adapt"
	"github.com/Goboolean/core-system.worker/internal/job/analyze"
	"github.com/Goboolean/core-system.worker/internal/job/fetch"
	modelExecute "github.com/Goboolean/core-system.worker/internal/job/model-execute"
	"github.com/Goboolean/core-system.worker/internal/job/transmit"
)

type Spec struct {
	FetchJobName     string
	ModelExecJobName string
	AdaptJobName     string
	ResAnalyzeJob    string
	TransmitJobName  string
}

func Build(spec Spec, UserParams *job.UserParams) (*Pipeline, error) {
	fetch, err := fetch.Create(spec.FetchJobName, UserParams)
	if err != nil {
		return nil, err
	}

	modelExec, err := modelExecute.Create(spec.ModelExecJobName, UserParams)
	if err != nil {
		return nil, err
	}

	adapt, err := adapt.Create(spec.AdaptJobName, UserParams)
	if err != nil {
		return nil, err
	}

	analyze, err := analyze.Create(spec.ResAnalyzeJob, UserParams)
	if err != nil {
		return nil, err
	}

	transmit, err := transmit.Create(spec.TransmitJobName, UserParams)
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
