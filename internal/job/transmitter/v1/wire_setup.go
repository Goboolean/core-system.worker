//go:build wireinject
// +build wireinject

package v1

import (
	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/infrastructure/influx"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	"github.com/google/wire"
)

func provideOrderEventDispatcher() (transmitter.OrderEventDispatcher, error) {
	return influx.NewOrderEventDispatcher(&influx.Opts{
		URL:        configuration.InfluxDBURL,
		Token:      configuration.InfluxDBToken,
		Org:        configuration.InfluxDBOrg,
		BucketName: configuration.InfluxDBOrderEventBucket,
	})
}

func provideAnnotationDispatcher() (transmitter.AnnotationDispatcher, error) {
	return influx.NewAnnotationDispatcher(&influx.Opts{
		URL:        configuration.InfluxDBURL,
		Token:      configuration.InfluxDBToken,
		Org:        configuration.InfluxDBOrg,
		BucketName: configuration.InfluxDBAnnotationBucket,
	})
}

func Create(p *job.UserParams) (transmitter.Transmitter, error) {
	wire.Build(wire.NewSet(
		provideOrderEventDispatcher,
		provideAnnotationDispatcher,
		NewCommon,
		wire.Bind(new(transmitter.Transmitter), new(*Common))))
	return &Common{}, nil
}
