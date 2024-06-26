package configuration

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ImportAppConfigFromFile(path string) (*AppConfig, error) {
	b, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("import app config: %w", err)
	}

	//act
	var config AppConfig
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, fmt.Errorf("import app config: %w", err)
	}

	return &config, nil
}

type AppConfig struct {
	Task           string         `yaml:"task"`
	TaskID         string         `yaml:"taskID"`
	InitialCapital int            `yaml:"initialCapital"`
	DataOrigin     DataOrigin     `yaml:"dataOrigin"`
	Model          ModelConfig    `yaml:"model"`
	Strategy       StrategyConfig `yaml:"strategy"`
}

type DataOrigin struct {
	TimeFrame      TimeFrame `yaml:"timeFrame"`
	ProductID      string    `yaml:"productID"`
	ProductType    string    `yaml:"productType"`
	StartTimestamp int64     `yaml:"startTimestamp"`
	EndTimestamp   int64     `yaml:"endTimestamp"`
}

type TimeFrame struct {
	Seconds int `yaml:"seconds"`
}

type ModelConfig struct {
	ID         string             `yaml:"ID"`
	BatchSize  int                `yaml:"batchSize"`
	OutputType string             `yaml:"outputType"`
	Params     map[string]float32 `yaml:"params"`
}

type StrategyConfig struct {
	ID        string             `yaml:"ID"`
	InputType string             `yaml:"inputType"`
	Params    map[string]float32 `yaml:"params"`
}
