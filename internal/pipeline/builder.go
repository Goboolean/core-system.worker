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
	"github.com/Goboolean/core-system.worker/internal/job/joiner"
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

func Build(config configuration.AppConfig) (Pipeline, error) {
	//step1: select pipeline
	t, err := selectPipeline(config)
	if err != nil {
		return nil, fmt.Errorf("build pipeline: %w", err)
	}

	switch t {
	case NormalPipeline:
		return buildNormal(config)
	case PipelineWithoutModel:
		return buildWithoutModel(config)
	default:
		return nil, ErrNotImplemented
	}
}

func selectPipeline(config configuration.AppConfig) (PipelineType, error) {
	if config.Model.ID == "" {
		return PipelineWithoutModel, nil
	}

	if config.Model.ID != "" {
		return NormalPipeline, nil
	}
	configStringBytes, err := yaml.Marshal(config)

	log.Error(fmt.Errorf("marshaling config: %w", err))
	return 0, fmt.Errorf("%w %s", ErrNoCompatiblePipeline, string(configStringBytes))
}

func buildNormal(config configuration.AppConfig) (*Normal, error) {
	p := extractUserParams(config)

	//job객체를 factory로부터 생성
	fetcher, err := fetcher.Create(extractFetcherSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	modelExecuter, err := executer.Create(extractModelExecterSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	joiner, err := joiner.NewBySequence(&p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	analyzer, err := analyzer.Create(extractAnalyzerSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	// transmitter 패키지는 factory가 없다. 그 이유는 transmit job은 한 가지 종류밖에 없기 때문이다.
	// 현재 생성자 미구현으로 dummy 객체로 대체
	transmitter, err := transmitter.NewFake()
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}

	isAdapterRequired := config.Model.OutputType != config.Strategy.InputType
	if isAdapterRequired {
		adapter, err := adapter.Create(adapter.Spec{
			InputType:  config.Model.OutputType,
			OutputType: config.Strategy.InputType,
		}, &p)
		if err != nil {
			return nil, fmt.Errorf("build normal pipeline: %w", err)
		}
		return newNormalWithAdapter(
			fetcher,
			joiner,
			modelExecuter,
			adapter,
			analyzer,
			transmitter,
		)
	} else {
		return newNormalWithoutAdapter(
			fetcher,
			joiner,
			modelExecuter,
			analyzer,
			transmitter,
		)
	}

}

func buildWithoutModel(config configuration.AppConfig) (*WithoutModel, error) {
	p := extractUserParams(config)

	//job객체를 factory로부터 생성
	fetcher, err := fetcher.Create(extractFetcherSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	analyzer, err := analyzer.Create(extractAnalyzerSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	// transmitter 패키지는 factory가 없다. 그 이유는 transmit job은 한 가지 종류밖에 없기 때문이다.
	// 현재 생성자 미구현으로 dummy 객체로 대체
	transmitter, err := transmitter.NewFake()
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}

	isAdapterRequired := config.DataOrigin.ProductType != config.Strategy.InputType
	if isAdapterRequired {
		adapter, err := adapter.Create(adapter.Spec{
			InputType:  config.DataOrigin.ProductType,
			OutputType: config.Strategy.InputType,
		}, &p)
		if err != nil {
			return nil, fmt.Errorf("build normal pipeline: %w", err)
		}
		return newWithoutModelWithAdapter(
			fetcher,
			adapter,
			analyzer,
			transmitter,
		)
	} else {
		return newWithoutModelWithoutAdapter(
			fetcher,
			analyzer,
			transmitter,
		)
	}

}

func extractFetcherSpec(config configuration.AppConfig) fetcher.Spec {

	var spec fetcher.Spec

	spec.ProductType = config.DataOrigin.ProductType
	spec.Task = config.Task
	return spec
}

func extractModelExecterSpec(config configuration.AppConfig) executer.Spec {

	var spec executer.Spec

	spec.OutputType = config.Model.OutputType
	return spec
}

func extractAnalyzerSpec(config configuration.AppConfig) analyzer.Spec {

	spec := analyzer.Spec{
		ID:        config.Strategy.ID,
		InputType: config.Strategy.InputType,
	}

	return spec
}

func extractUserParams(config configuration.AppConfig) job.UserParams {

	var p = job.UserParams{
		job.StartDate: fmt.Sprint(config.DataOrigin.StartTimestamp),
		job.EndDate:   fmt.Sprint(config.DataOrigin.EndTimestamp),
		job.BatchSize: fmt.Sprint(config.Model.BatchSize),
		job.ProductID: config.DataOrigin.ProductID,
	}

	for k, v := range config.Model.Params {
		p[strings.Join([]string{"model", k}, ".")] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	}

	for k, v := range config.Strategy.Params {
		p[strings.Join([]string{"strategy", k}, ".")] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	}

	return p
}
