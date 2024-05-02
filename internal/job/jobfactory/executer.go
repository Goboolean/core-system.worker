package jobfactory

import (
	"fmt"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
)

type executerProvider func(p *job.UserParams) (executer.ModelExecutor, error)

// wire의 한계로 여기서 수동DI 합니다.
var executerProviderRepo = map[executer.Spec]executerProvider{
	executer.Spec{OutputType: "candlestick"}: func(p *job.UserParams) (executer.ModelExecutor, error) {

		//더미
		kServeConfig := &resolver.ConfigMap{}
		kServeClient, err := infrastructure.NewKServeClient(kServeConfig)
		if err != nil {
			return nil, err
		}
		return executer.NewMock(kServeClient, *p)
	},
	executer.Spec{OutputType: "value"}: func(p *job.UserParams) (executer.ModelExecutor, error) {
		return executer.Dummy{}, nil
	},
	executer.Spec{OutputType: "proveDist"}: func(p *job.UserParams) (executer.ModelExecutor, error) {
		return executer.Dummy{}, nil
	},
}

func CreateExecuter(spec executer.Spec, p *job.UserParams) (executer.ModelExecutor, error) {

	var provider, ok = executerProviderRepo[spec]
	if !ok {
		return nil, NotFoundJob
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create model execute job: %w", err)
	}

	return f, nil
}
