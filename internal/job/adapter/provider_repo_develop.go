//go:build develop

package adapter

import "github.com/Goboolean/core-system.worker/internal/job"

var providerRepo = map[Spec]jobProvider{
	{InputType: "candlestick", OutputType: "candlestick"}: func(p *job.UserParams) (Adapter, error) {
		return Dummy{}, nil
	},
}
