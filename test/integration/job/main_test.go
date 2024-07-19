package job

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	url         = ""
	token       = configuration.InfluxDBToken
	org         = configuration.InfluxDBOrg
	tradeBucket = configuration.InfluxDBTradeBucket
)

var influxC *InfluxContainer

type InfluxContainer struct {
	testcontainers.Container
}

func InitializeInfluxContainer(ctx context.Context) (*InfluxContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "influxdb:latest",
		ExposedPorts: []string{"8086"},
		WaitingFor:   wait.ForListeningPort("8086"),
		Env: map[string]string{
			"DOCKER_INFLUXDB_INIT_MODE":        "setup",
			"DOCKER_INFLUXDB_INIT_USERNAME":    "admin",
			"DOCKER_INFLUXDB_INIT_PASSWORD":    "password",
			"DOCKER_INFLUXDB_INIT_ORG":         org,
			"DOCKER_INFLUXDB_INIT_BUCKET":      "bucket",
			"DOCKER_INFLUXDB_INIT_ADMIN_TOKEN": token,
			"INFLUXD_LOG_LEVEL":                "debug",
		},
	}

	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	cmds := []string{
		fmt.Sprintf(`influx bucket create --token %s --org %s --name %s`, token, org, tradeBucket),
	}

	for _, c := range cmds {
		_, res, err := container.Exec(ctx, []string{"sh", "-c", c})
		if err != nil {
			return nil, fmt.Errorf("create bucket: fail to create bucket")
		}

		str, err := io.ReadAll(res)
		if err != nil {
			return nil, fmt.Errorf("read result: fail to read out string")
		}

		fmt.Println(string(str))
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "8086")
	if err != nil {
		return nil, err
	}
	fmt.Println(port)

	url = fmt.Sprintf("http://%s:%s", ip, port.Port())

	return &InfluxContainer{Container: container}, nil
}

const (
	testStockID   = "stock.aapl.usa"
	testTimeFrame = "1m"
)

func TestMain(m *testing.M) {
	var err error
	influxC, err = InitializeInfluxContainer(context.Background())
	if err != nil {
		panic(err)
	}
	m.Run()
	err = influxC.Terminate(context.Background())
	if err != nil {
		panic(err)
	}
}
