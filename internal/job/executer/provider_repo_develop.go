//go:build develop

package executer

import "github.com/Goboolean/core-system.worker/internal/job"

var providerRepo = map[Spec]jobProvider{
	{OutputType: "candlestick"}: initalizeMock,
	{OutputType: "stub"}: func(p *job.UserParams) (ModelExecutor, error) {
		return NewStub(p)
	},
}
