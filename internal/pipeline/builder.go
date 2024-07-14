package pipeline

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/adapter"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/joiner"
	v1 "github.com/Goboolean/core-system.worker/internal/job/transmitter/v1"
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

// Build constructs the appropriate pipeline based on the application config.
//
// The build process proceeds as follows:
//
// 1. Select Pipeline
//
// Based on the configuration, determine which pipeline the user desires by verifying omitted or additional settings.
// For example, if there is no configuration for a model, it indicates that the user wants a pipeline without a model.
//
// 2. Create Job
//
// Extract specs that distinguishes jobs for each stage from the application configuration,
// and pass the spec and other parameters to Create to generate the necessary jobs for the selected pipeline.
//
// 3. Create Pipeline
//
// Inject the generated jobs into the appropriate pipeline to create the pipeline.
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

// selectPipeline determine which pipeline the user desires by verifying omitted or additional settings.
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

// buildNormal builds normal pipeline
func buildNormal(config configuration.AppConfig) (*Normal, error) {
	p := extractUserParams(config)

	//job객체를 factory로부터 생성
	fetcher, err := fetcher.Create(extractFetcherSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	modelExecuter, err := executer.Create(extractModelExecuterSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	joiner, err := joiner.NewByTime(&p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}
	analyzer, err := analyzer.Create(extractAnalyzerSpec(config), &p)
	if err != nil {
		return nil, fmt.Errorf("build normal pipeline: %w", err)
	}

	transmitter, err := v1.Create(&p)
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
		return NewNormalWithAdapter(
			fetcher,
			joiner,
			modelExecuter,
			adapter,
			analyzer,
			transmitter,
		)
	} else {
		return NewNormalWithoutAdapter(
			fetcher,
			joiner,
			modelExecuter,
			analyzer,
			transmitter,
		)
	}

}

// buildWithoutModel builds pipeline without model
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

	transmitter, err := v1.Create(&p)
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
		return NewWithoutModelWithAdapter(
			fetcher,
			adapter,
			analyzer,
			transmitter,
		)
	} else {
		return NewWithoutModelWithoutAdapter(
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

func extractModelExecuterSpec(config configuration.AppConfig) executer.Spec {

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
		job.TaskID:    config.TaskID,
	}

	for k, v := range config.Model.Params {
		p[strings.Join([]string{"model", k}, ".")] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	}

	for k, v := range config.Strategy.Params {
		p[strings.Join([]string{"strategy", k}, ".")] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	}

	timeFrame := (time.Duration(config.DataOrigin.TimeFrame.Seconds) * time.Second).String()
	timeFrame = strings.Replace(timeFrame, "m0s", "m", 1)
	timeFrame = strings.Replace(timeFrame, "h0m", "h", 1)

	p[job.TimeFrame] = fmt.Sprint(timeFrame)
	return p
}
