package kserve

// InferenceRes는 kserve로부터 받은 응답을 담기 위한 구조체입니다.
type InferenceRes struct {
	ModelName    string           `json:"model_name"`
	ModelVersion string           `json:"model_version,omitempty"`
	Id           string           `json:"id"`
	Parameters   string           `json:"parameters,omitempty"`
	Outputs      []ResponseOutput `json:"outputs"`
}

// ResponseOutput는 kserve로부터 받은 응답의 출력을 담기 위한 구조체입니다.
type ResponseOutput struct {
	Name       string      `json:"name"`
	Shape      []int       `json:"shape"`
	DataType   string      `json:"datatype"`
	Parameters string      `json:"parameters,omitempty"`
	Data       interface{} `json:"data"`
}
