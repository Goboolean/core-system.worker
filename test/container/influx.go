package container

import (
	"context"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type InfluxContainer struct {
	testcontainers.Container
	URL string
}

var influxEnvs = map[string]string{
	"DOCKER_INFLUXDB_INIT_MODE":        "setup",
	"DOCKER_INFLUXDB_INIT_USERNAME":    "admin",
	"DOCKER_INFLUXDB_INIT_PASSWORD":    "password",
	"DOCKER_INFLUXDB_INIT_ORG":         os.Getenv("INFLUXDB_ORG"),
	"DOCKER_INFLUXDB_INIT_BUCKET":      "bucket",
	"DOCKER_INFLUXDB_INIT_ADMIN_TOKEN": os.Getenv("INFLUXDB_TOKEN"),
	"INFLUXD_LOG_LEVEL":                "debug",
}

func InitInfluxContainer(ctx context.Context, buckets ...string) (*InfluxContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "influxdb:latest",
		ExposedPorts: []string{"8086"},
		WaitingFor:   wait.ForListeningPort("8086"),
		Env:          influxEnvs,
	}

	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return nil, fmt.Errorf("create container: fail to create container %v", err)
	}

	for _, b := range buckets {
		err := createBucket(ctx, container, b)
		if err != nil {
			return nil, fmt.Errorf("create container: fail to create container %v", err)
		}
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("create container: fail to create container %v", err)
	}

	port, err := container.MappedPort(ctx, "8086")
	if err != nil {
		return nil, fmt.Errorf("create container: fail to create container %v", err)
	}

	log.WithField("port", port.Port()).Debug("Influx container is created")

	url := fmt.Sprintf("http://%s:%s", ip, port.Port())

	return &InfluxContainer{Container: container, URL: url}, nil
}

func createBucket(ctx context.Context, container testcontainers.Container, b string) error {

	log.WithField("bucketName", b).Debug("Creating Influxdb bucket")

	cmd := fmt.Sprintf(`influx bucket create --token %s --org %s --name %s`, os.Getenv("INFLUXDB_TOKEN"), os.Getenv("INFLUXDB_ORG"), b)

	_, res, err := container.Exec(ctx, []string{"sh", "-c", cmd})
	if err != nil {
		return fmt.Errorf("create bucket: fail to create bucket")
	}

	str, err := io.ReadAll(res)
	if err != nil {
		return fmt.Errorf("read result: fail to read out string")
	}

	println(string(str))

	return nil
}
