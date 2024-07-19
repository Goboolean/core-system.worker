package container

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
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
	"DOCKER_INFLUXDB_INIT_ORG":         configuration.InfluxDBOrg,
	"DOCKER_INFLUXDB_INIT_BUCKET":      "bucket",
	"DOCKER_INFLUXDB_INIT_ADMIN_TOKEN": configuration.InfluxDBToken,
	"INFLUXD_LOG_LEVEL":                "debug",
}

func InitInfluxContainerWithRandomPort(ctx context.Context, buckets ...string) (*InfluxContainer, error) {
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

func InitInfluxContainerWithPortBinding(ctx context.Context, port string, buckets ...string) (*InfluxContainer, error) {

	portStringWithProtocol := strings.Join([]string{port, "tcp"}, "/")
	req := testcontainers.ContainerRequest{
		Image:        "influxdb:latest",
		ExposedPorts: []string{"8086"},
		WaitingFor:   wait.ForListeningPort("8086"),
		Env:          influxEnvs,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				nat.Port(portStringWithProtocol): {{HostIP: "0.0.0.0", HostPort: portStringWithProtocol}},
			}
		},
	}

	waitForPortToBeClosed(ctx, port)

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

	fmt.Println(port)

	url := fmt.Sprintf("http://%s:%s", ip, port)

	return &InfluxContainer{Container: container, URL: url}, nil
}

func createBucket(ctx context.Context, container testcontainers.Container, b string) error {

	log.WithField("bucketName", b).Debug("Creating Influxdb bucket")

	cmd := fmt.Sprintf(`influx bucket create --token %s --org %s --name %s`, configuration.InfluxDBToken, configuration.InfluxDBOrg, b)

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

func waitForPortToBeClosed(ctx context.Context, port string) {
	for {
		exit := false
		select {
		case <-ctx.Done():
			return
		default:
			if !isPortOpen(port) {
				exit = true
			}
			time.Sleep(100 * time.Millisecond)
		}
		if exit {
			break
		}
	}
}

func isPortOpen(port string) bool {
	address := strings.Join([]string{"localhost", port}, ":")
	// Listen이 가능하다 == 포트가 점유되지 않았다.
	conn, err := net.Listen("tcp", address)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}
