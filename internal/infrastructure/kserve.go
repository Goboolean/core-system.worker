package infrastructure

import "github.com/Goboolean/common/pkg/resolver"

type KServeClientImpl struct {
}

// KServeClient는 KServe에 요청을 보내고 받는 역할을 하는 인터페이스입니다.
type KServeClient interface {
	SetModelName(name string)
	RequestInference(ctx context.Context, shape []int, input []float32) (output []float32, err error)
}

func NewKServeClient(c *resolver.ConfigMap) *KServeClientImpl {
	return &KServeClientImpl{}
}
