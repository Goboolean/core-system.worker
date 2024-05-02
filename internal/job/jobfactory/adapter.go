package jobfactory

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
	adapter "github.com/Goboolean/core-system.worker/internal/job/analyzer"
)

type adapterProvider func(p *job.UserParams) (adapter.Analyzer, error)

// wire의 한계로 여기서 수동DI 합니다.
var adapterProviderRepo = map[adapter.Spec]adapterProvider{
	adapter.Spec{Id: "abc", InputType: "candlestick"}: func(p *job.UserParams) (adapter.Analyzer, error) {
		return adapter.Dummy{}, nil
	},
}

func CreateAdapter(spec adapter.Spec, p *job.UserParams) (adapter.Analyzer, error) {

	var provider, ok = adapterProviderRepo[spec]
	if !ok {
		return nil, NotFoundJob
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create analyze job: %w", err)
	}

	return f, nil
}
