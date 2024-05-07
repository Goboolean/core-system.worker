//go:build wireinject
// +build wireinject

package executer

import (
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/infrastructure/kserve"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/google/wire"
)

type kServeConfig resolver.ConfigMap

func provideKServeConfig() kServeConfig {
	return kServeConfig(resolver.ConfigMap{})
}

func provideKServe(c kServeConfig) (*kserve.ClientImpl, error) {
	in := resolver.ConfigMap(c)
	return kserve.NewClient(&in)
}

func initalizeMock(p *job.UserParams) (ModelExecutor, error) {
	wire.Build(
		provideKServeConfig,
		provideKServe,
		NewMock,
		wire.Bind(new(ModelExecutor), new(*Mock)),
		wire.Bind(new(kserve.Client), new(*kserve.ClientImpl)),
	)
	return nil, nil
}
