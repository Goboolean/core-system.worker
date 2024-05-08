package kserve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/dto"
)

var defaultKSserveClient = &ClientImpl{
	host:   "",
	param1: 0.0,
	param2: 1.0,
	http: &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 600 * time.Second,
		},
	},
}

// Client is an interface that defines the role of sending and receiving requests to KServe.
type Client interface {
	SetModelName(name string)
	RequestInference(ctx context.Context, shape []int, input []float32) (output []float32, err error)
}

// ClientImpl is a struct that represents the implementation of the KServeClient interface.
type ClientImpl struct {
	modelId           string
	host              string
	inferenceEndpoint *url.URL

	param1 float32
	param2 float32

	http *http.Client
}

// NewClient creates a new instance of KServeClientImpl.
func NewClient(c *resolver.ConfigMap) (*ClientImpl, error) {

	//default value
	instance := defaultKSserveClient

	id, err := c.GetStringKey("modelID")
	if err != nil {
		return nil, fmt.Errorf("create kserve client: %w", err)
	}

	host, err := c.GetStringKey("host")
	if err != nil {
		return nil, fmt.Errorf("create kserve client: %w", err)
	}

	param1, err := c.GetFloatKey("param1")
	if err != nil {
		return nil, fmt.Errorf("create kserve client: %w", err)
	}

	param2, err := c.GetFloatKey("param2")
	if err != nil {
		return nil, fmt.Errorf("create kserve client: %w", err)
	}

	instance.host = host
	instance.modelId = id
	instance.param1 = float32(param1)
	instance.param2 = float32(param2)

	return instance, nil
}

// SetModelName sets the model name for the KServeClientImpl instance.
func (c *ClientImpl) SetModelName(name string) {
	c.inferenceEndpoint = generateInferenceUrl(c.host, name)
}

// RequestInference sends an inference request to KServe and returns the output and any error that occurred.
func (c *ClientImpl) RequestInference(ctx context.Context, shape []int, input []float32) ([]float32, error) {
	bodyJson, err := json.Marshal(
		dto.InferenceReq{
			Name:     "input",
			Shape:    shape,
			DataType: "FP32",
			Data:     input,
		})

	if err != nil {
		return nil, fmt.Errorf("inference: falid to create request body  %w", err)
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    c.inferenceEndpoint,
		Body:   io.NopCloser(bytes.NewBuffer(bodyJson)),
	}
	req = req.WithContext(ctx)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("inference: %w", err)
	}

	var b []byte
	var out dto.InferenceRes
	if _, err := res.Body.Read(b); err != nil {
		return nil, fmt.Errorf("inference: falid to read responce body %w", err)
	}

	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("inference: falid to unmarshal responce %w", err)

	}

	return out.Outputs[0].Data.([]float32), nil
}

func generateInferenceUrl(host, modelName string) *url.URL {
	return &url.URL{
		Host: host,
		Path: fmt.Sprintf("v2/models/%s/infer", modelName),
	}
}
