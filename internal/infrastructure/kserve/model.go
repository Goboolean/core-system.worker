package kserve

// InferenceReq is a struct for sending requests to kserve.
type InferenceReq struct {
	Name       string      `json:"string"`
	Shape      []int       `json:"shape"`
	DataType   string      `json:"datatype"`
	Parameters string      `json:"parameters,omitempty"`
	Data       interface{} `json:"data"`
}

// InferenceRes is a struct for storing responses received from kserve.
type InferenceRes struct {
	ModelName    string           `json:"model_name"`
	ModelVersion string           `json:"model_version,omitempty"`
	ID           string           `json:"ID"`
	Parameters   string           `json:"parameters,omitempty"`
	Outputs      []ResponseOutput `json:"outputs"`
}

// ResponseOutput is a struct for storing the output of responses received from kserve.
type ResponseOutput struct {
	Name       string      `json:"name"`
	Shape      []int       `json:"shape"`
	DataType   string      `json:"datatype"`
	Parameters string      `json:"parameters,omitempty"`
	Data       interface{} `json:"data"`
}
