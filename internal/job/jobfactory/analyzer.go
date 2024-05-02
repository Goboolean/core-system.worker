package jobfactory

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
)

type analyzerProvider func(p *job.UserParams) (analyzer.Analyzer, error)

// wire의 한계로 여기서 수동DI 합니다.
var analyzerProviderRepo = map[analyzer.Spec]adapterProvider{
	analyzer.Spec{Id: "abc", InputType: "candlestick"}: func(p *job.UserParams) (analyzer.Analyzer, error) {
		return analyzer.Dummy{}, nil
	},
}

func CreateAnalyzer(spec analyzer.Spec, p *job.UserParams) (analyzer.Analyzer, error) {

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
