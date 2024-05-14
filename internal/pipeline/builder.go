package pipeline

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/adapter"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type PipelineType int

const (
	NormalPipeline PipelineType = iota + 1
	PipelineWithoutModel
)

var (
	ErrNoCompatiblePipeline = errors.New("select pipeline: there are no compatible pipeline")
	ErrNotImplemented       = errors.New("select pipeline: selected pipeline is not implemented")
)

func Build(config configuration.AppConfig) (*Pipeline, error) {
	//step1: select pipeline
	t, err := selectPipeline(config)
	if err != nil {
		return nil, fmt.Errorf("build pipeline: %w", err)
	}

	switch t {
	case NormalPipeline:
		return buildNormal(config)
	default:
		return nil, ErrNotImplemented
	}
}

func selectPipeline(config configuration.AppConfig) (PipelineType, error) {
	if config.Model.Id == "" {
		return PipelineWithoutModel, nil
	}
	if config.Model.Id != "" {
		return NormalPipeline, nil
	}
	configStringBytes, err := yaml.Marshal(config)

	log.Error(fmt.Errorf("marshaling config: %w", err))
	return 0, fmt.Errorf("%w %s", ErrNoCompatiblePipeline, string(configStringBytes))
}

func buildNormal(config configuration.AppConfig) (*Pipeline, error) {
	p := extractUserParams(config)

	//job객체를 factory로부터 생성
	fetchJob, err := fetcher.Create(extractFetchSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	modelExecJob, err := executer.Create(extractModelExecSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	isAdapterRequred, adapterSpec := extractAdapterSpec(config)
	adpatJob, err := adapter.Create(adapterSpec, &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	analyzeJob, err := analyzer.Create(extractAnalyzerSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	// transmitter 패키지는 factory가 없다. 그 이유는 transmit job은 한 가지 종류밖에 없기 때문이다.
	// 현재 생성자 미구현으로 dummy 객체로 대체
	transmitJob := transmitter.Dummy{}

	if isAdapterRequred {
		if err != nil {
			return nil, fmt.Errorf("build normal pipeline: %w", err)
		}
		return newNormalWithAdapter(
			fetchJob,
			modelExecJob,
			adpatJob,
			analyzeJob,
			transmitJob,
		)
	} else {
		return newNormalWithoutAdapter(
			fetchJob,
			modelExecJob,
			analyzeJob,
			transmitJob,
		)
	}

}

func extractFetchSpec(config configuration.AppConfig) fetcher.Spec {

	var spec fetcher.Spec

	spec.ProductType = config.DataOrigin.ProductType
	spec.Task = spec.Task
	return spec
}

func extractModelExecSpec(config configuration.AppConfig) executer.Spec {

	var spec executer.Spec

	spec.OutputType = config.Model.OutputType
	return spec
}

func extractAdapterSpec(config configuration.AppConfig) (bool, adapter.Spec) {

	isRequred := config.Model.OutputType == config.Strategy.InputType
	sepc := adapter.Spec{
		InputType:  config.Model.OutputType,
		OutputType: config.Strategy.InputType,
	}

	return isRequred, sepc

}

func extractAnalyzerSpec(config configuration.AppConfig) analyzer.Spec {

	spec := analyzer.Spec{
		Id:        config.Strategy.Id,
		InputType: config.Strategy.InputType,
	}

	return spec
}

func extractUserParams(config configuration.AppConfig) job.UserParams {

	var p = job.UserParams{
		"startDate": string(config.DataOrigin.StartTimestamp),
		"endDate":   string(config.DataOrigin.EndTimestamp),
		"batchSize": string(config.Model.BatchSize),
		"productId": config.DataOrigin.ProductId,
	}

	for k, v := range config.Model.Params {
		p[strings.Join([]string{"model", k}, ".")] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	}

	for k, v := range config.Strategy.Params {
		p[strings.Join([]string{"stretage", k}, ".")] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	}

	return p
}
