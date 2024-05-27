//go:build develop

package analyzer

import "github.com/Goboolean/core-system.worker/internal/job"

var providerRepo = map[Spec]jobProvider{
	{ID: "stub", InputType: "candlestick"}: func(p *job.UserParams) (Analyzer, error) {
		return NewStub(p)
	},
}
