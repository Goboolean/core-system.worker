package configuration

type AppConfig struct {
	Task       string         `yaml:"task"`
	DataOrigin DataOrigin     `yaml:"dataOrigin"`
	Model      ModelConfig    `yaml:"model"`
	Strategy   StrategyConfig `yaml:"strategy"`
}

type DataOrigin struct {
	TimeFrame      TimeFrame `yaml:"timeFrame"`
	ProductType    string    `yaml:"productType"`
	StartTimestamp int64     `yaml:"startTimestamp"`
	EndTimestamp   int64     `yaml:"endTimestamp"`
}

type TimeFrame struct {
	Seconds int `yaml:"seconds"`
}

type ModelConfig struct {
	Id         string             `yaml:"id"`
	BatchSize  int                `yaml:"batchSize"`
	OutputType string             `yaml:"outputType"`
	Params     map[string]float32 `yaml:"params"`
}

type StrategyConfig struct {
	Id        string             `yaml:"id"`
	InputType string             `yaml:"inputType"`
	Params    map[string]float32 `yaml:"params"`
}
