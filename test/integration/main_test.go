package integration

import (
	"context"
	"os"
	"testing"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stretchr/testify/assert"
)

var rawInfluxClient influxdb2.Client

var (
	influxDBUrl   = os.Getenv("INFLUXDB_URL")
	influxDBOrg   = os.Getenv("INFLUXDB_ORG")
	influxDBToken = os.Getenv("INFLUXDB_TOKEN")

	tradeBucket      = os.Getenv("INFLUXDB_TRADE_BUCKET")
	orderBucket      = os.Getenv("INFLUXDB_ORDER_EVENT_BUCKET")
	annotationBucket = os.Getenv("INFLUXDB_ANNOTATION_EVENT_BUCKET")
)

func RecreateBucket(client influxdb2.Client, orgName, bucketName string) error {
	org, err := client.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	bucket, err := client.BucketsAPI().FindBucketByName(context.Background(), bucketName)
	if err != nil {
		return nil
	}

	client.BucketsAPI().DeleteBucket(context.Background(), bucket)
	_, err = client.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)

	return err
}

func TestPing(t *testing.T) {
	ok, err := rawInfluxClient.Ping(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestMain(m *testing.M) {
	rawInfluxClient = influxdb2.NewClient(influxDBUrl, influxDBToken)
	m.Run()
	rawInfluxClient.Close()
}
