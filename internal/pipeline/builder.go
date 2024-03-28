package pipeline

import "github.com/Goboolean/core-system.worker/internal/job"

type Spec struct {
	FetchJobName     string
	ModelExecJobName string
	AdaptJobName     string
	ResAnalyzeJob    string
	TransmitJobName  string
}

type PipelineBuilder struct {
	factory job.JobFactory
}

func NewPipelineBuilder(jobFactory job.JobFactory) *PipelineBuilder {
	return &PipelineBuilder{
		factory: jobFactory,
	}
}

func (b *PipelineBuilder) Build(spec Spec, UserParams *job.UserParams) (*Pipeline, error) {
	fetch, err := b.factory.CreateJob(spec.FetchJobName, UserParams)
	if err != nil {
		return nil, err
	}

	modelExec, err := b.factory.CreateJob(spec.ModelExecJobName, UserParams)
	if err != nil {
		return nil, err
	}

	adapt, err := b.factory.CreateJob(spec.AdaptJobName, UserParams)
	if err != nil {
		return nil, err
	}

	analyze, err := b.factory.CreateJob(spec.ResAnalyzeJob, UserParams)
	if err != nil {
		return nil, err
	}

	transmit, err := b.factory.CreateJob(spec.TransmitJobName, UserParams)
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
