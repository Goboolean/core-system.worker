package influx_test

import (
	"context"
	"os"
	"testing"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	url    = os.Getenv("INFLUXDB_URL")
	bucket = os.Getenv("INFLUXDB_BUCKET")
	token  = os.Getenv("INFLUXDB_TOKEN")
	org    = os.Getenv("INFLUXDB_ORG")
)

var rawInfluxDBClient influxdb2.Client

var (
	productID = "stock.aapl.usa"
	taskID    = "2024-07-05-test"
)

func TestMain(m *testing.M) {
	rawInfluxDBClient = influxdb2.NewClient(url, token)
	m.Run()
	rawInfluxDBClient.Close()
}

func RecreateBucket(client influxdb2.Client, orgName, bucketName string) error {

	org, err := client.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	bucket, err := client.BucketsAPI().FindBucketByName(context.Background(), bucketName)
	if err != nil {
		return nil
	}

	if err := client.BucketsAPI().DeleteBucket(context.Background(), bucket); err != nil {
		return err
	}

	_, err = client.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)

	return err
}
