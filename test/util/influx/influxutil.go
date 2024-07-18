package influx

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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

	if err := client.BucketsAPI().DeleteBucket(context.Background(), bucket); err != nil {
		return err
	}
	_, err = client.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)

	return err
}

func CountRecordsInMeasurement(client influxdb2.Client, orgName, bucketName, measurement string) (int, error) {

	q, err := client.QueryAPI(orgName).
		Query(context.Background(),
			fmt.Sprintf(
				`from(bucket: "%s")
				|> range(start:0)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> count()`, bucketName, measurement))
	if err != nil {
		return 0, err
	}

	num := int64(0)

	// 각 record 별 count에서 최댓값을 찾는다.
	for q.Next() {
		fmt.Println("CountRecordsInMeasurement: records: ", q.Record().Values())
		val := q.Record().ValueByKey("_value").(int64)
		if val > num {
			num = val
		}
	}

	return int(num), nil
}
