package kserve

// InferenceReq는 kserve에 요청을 보내기 위한 구조체입니다.
type InferenceReq struct {
	Name       string        `json:"string"`
	Shape      []int         `json:"shape"`
	DataType   string        `json:"datatype"`
	Parameters string        `json:"parameters,omitempty"`
	Data       []interface{} `json:"data"`
}
