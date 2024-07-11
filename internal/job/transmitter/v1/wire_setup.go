//go:build wireinject
// +build wireinject

package v1

import (
	"os"

	"github.com/Goboolean/core-system.worker/internal/infrastructure/influx"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	"github.com/google/wire"
)

func provideOrderEventDispatcher() (transmitter.OrderEventDispatcher, error) {
	return influx.NewOrderEventDispatcher(&influx.Opts{
		RUL:        os.Getenv("INFLUXDB_URL"),
		Token:      os.Getenv("INFLUXDB_TOKEN"),
		Org:        os.Getenv("INFLUXDB_ORG"),
		BucketName: os.Getenv("INFLUXDB_ORDER_EVENT_BUCKET"),
	})
}

func provideAnnotationDispatcher() (transmitter.AnnotationDispatcher, error) {
	return influx.NewAnnotationDispatcher(&influx.Opts{
		RUL:        os.Getenv("INFLUXDB_URL"),
		Token:      os.Getenv("INFLUXDB_TOKEN"),
		Org:        os.Getenv("INFLUXDB_ORG"),
		BucketName: os.Getenv("INFLUXDB_ANNOTATION_BUCKET"),
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
