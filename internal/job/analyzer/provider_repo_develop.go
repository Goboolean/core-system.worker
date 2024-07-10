//go:build develop

package analyzer

import "github.com/Goboolean/core-system.worker/internal/job"

var providerRepo = map[Spec]jobProvider{
	{ID: "stub", InputType: "stub"}: func(p *job.UserParams) (Analyzer, error) {
		return NewStub(p)
	},
	{ID: "stub", InputType: "stockStub"}: func(p *job.UserParams) (Analyzer, error) {
		return NewStub(p)
	},
	{ID: "stub", InputType: "stock"}: func(p *job.UserParams) (Analyzer, error) {
		return NewStub(p)
	},
}
