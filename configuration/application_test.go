package configuration_test

import (
	"os"
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestUnmarshal(t *testing.T) {
	//arrange
	b, err := os.ReadFile("../config.example.yml")

	if err != nil {
		t.Error(err)
	}

	//act
	var AppConfig configuration.AppConfig
	err = yaml.Unmarshal(b, &AppConfig)

	//assert
	assert.NoError(t, err)
	assert.Equal(t, "backTest", AppConfig.Task)
	assert.Equal(t, configuration.DataOrigin{
		TimeFrame:      configuration.TimeFrame{Seconds: 1},
		ProductType:    "stock",
		StartTimestamp: 12345678,
		EndTimestamp:   12345678,
	}, AppConfig.DataOrigin)
	assert.Equal(t, configuration.ModelConfig{
		Id:         "goooo",
		BatchSize:  100,
		OutputType: "candlestick",
		Params: map[string]float32{
			"param1": 3.14,
		},
	}, AppConfig.Model)
	assert.Equal(t, configuration.StrategyConfig{
		Id:        "boolean",
		InputType: "candlestick",
		Params: map[string]float32{
			"param1": 3.14,
		},
	}, AppConfig.Strategy)
}
