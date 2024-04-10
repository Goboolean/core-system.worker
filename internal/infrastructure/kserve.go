package infrastructure

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
	kserve "github.com/Goboolean/core-system.worker/internal/dto/Kserve"
)

// KServeClient is an interface that defines the role of sending and receiving requests to KServe.
type KServeClient interface {
	SetModelName(name string)
	RequestInference(shape []int, input []float32) (output []float32, err error)
}

// KServeClientImpl is a struct that represents the implementation of the KServeClient interface.
type KServeClientImpl struct {
	modelId           string
	host              string
	inferenceEndpoint *url.URL

	param1 float32
	param2 float32

	http *http.Client
}

// NewKServeClient creates a new instance of KServeClientImpl.
func NewKServeClient(c *resolver.ConfigMap) (*KServeClientImpl, error) {

	//default value
	instance := &KServeClientImpl{
		host:   "",
		param1: 0.0,
		param2: 1.0,
		http: &http.Client{
			Transport: &http.Transport{
				IdleConnTimeout: 600 * time.Second,
			},
		},
	}

	host, err := c.GetStringKey("host")
	if err != nil {
		return nil, err
	}

	param1, err := c.GetFloatKey("param1")
	if err != nil {
		return nil, err
	}

	param2, err := c.GetFloatKey("param2")
	if err != nil {
		return nil, err
	}

	instance.host = host
	instance.param1 = float32(param1)
	instance.param2 = float32(param2)

	return instance, nil
}

// SetModelName sets the model name for the KServeClientImpl instance.
func (c *KServeClientImpl) SetModelName(name string) {
	c.inferenceEndpoint = generateInferenceUrl(c.host, name)
}

// RequestInference sends an inference request to KServe and returns the output and any error that occurred.
func (c *KServeClientImpl) RequestInference(ctx context.Context, shape []int, input []float32) (output []float32, err error) {
	bodyJson, err := json.Marshal(
		kserve.InferenceReq{
			Name:     "input",
			Shape:    shape,
			DataType: "FP32",
			Data:     input,
		})

	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    c.inferenceEndpoint,
		Body:   io.NopCloser(bytes.NewBuffer(bodyJson)),
	}
	req = req.WithContext(ctx)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	var b []byte
	var out kserve.InferenceRes
	res.Body.Read(b)

	json.Unmarshal(b, &out)
	return out.Outputs[0].Data.([]float32), nil
}

func generateInferenceUrl(host, modelName string) *url.URL {
	return &url.URL{
		Host: host,
		Path: fmt.Sprint("v2/models/%s/infer", modelName),
	}
}
