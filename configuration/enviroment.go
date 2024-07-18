package configuration

import "os"

var (
	InfluxDBURL              = os.Getenv("INFLUXDB_URL")
	InfluxDBToken            = os.Getenv("INFLUXDB_TOKEN")
	InfluxDBOrg              = os.Getenv("INFLUXDB_ORG")
	InfluxDBTradeBucket      = os.Getenv("INFLUXDB_BUCKET")
	InfluxDBOrderEventBucket = os.Getenv("INFLUXDB_ORDER_EVENT_BUCKET")
	InfluxDBAnnotationBucket = os.Getenv("INFLUXDB_ANNOTATION_BUCKET")
)
