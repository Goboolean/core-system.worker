//go:build wireinject
// +build wireinject

package executer

import (
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/google/wire"
)

type kServeConfig resolver.ConfigMap

func provideKServeConfig() kServeConfig {
	return kServeConfig(resolver.ConfigMap{})
}

func provideKServe(c kServeConfig) (*infrastructure.KServeClientImpl, error) {
	in := resolver.ConfigMap(c)
	return infrastructure.NewKServeClient(&in)
}

func initalizeMock(p *job.UserParams) (ModelExecutor, error) {
	wire.Build(
		provideKServeConfig,
		provideKServe,
		NewMock,
		wire.Bind(new(ModelExecutor), new(*Mock)),
		wire.Bind(new(infrastructure.KServeClient), new(*infrastructure.KServeClientImpl)),
	)
	return nil, nil
}
