package influx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/fetch-system.IaC/pkg/influx/mapper"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write" // Change the import statement to use the correct version of the Point type
)

var ErrBucketNotExist = errors.New("create dispatcher: influxdb bucket does not exists")

type Opts struct {
	URL        string
	Token      string
	Org        string
	BucketName string
}

type OrderEventDispatcher struct {
	// Dispatch dispatches the given order event.
	client influxdb2.Client
	writer api.WriteAPI
}

func NewOrderEventDispatcher(o *Opts) (*OrderEventDispatcher, error) {

	if o.URL == "" {
		return nil, fmt.Errorf("create influx db client: Required field Url is blank")
	}

	if o.Token == "" {
		return nil, fmt.Errorf("create influx db client: Required field Token is blank")
	}

	if o.Org == "" {
		return nil, fmt.Errorf("create influx db client: Required field Url is blank")
	}

	if o.BucketName == "" {
		return nil, fmt.Errorf("create influx db client: Required field TradeBucketName is blank")
	}

	client := influxdb2.NewClient(o.URL, o.Token)

	instance := &OrderEventDispatcher{
		client: client,
		writer: client.WriteAPI(o.Org, o.BucketName),
	}

	if !instance.bucketExists(o.BucketName) {
		return nil, ErrBucketNotExist
	}

	return instance, nil
}

func (d *OrderEventDispatcher) Dispatch(taskID string, event *model.OrderEvent) {
	d.writer.WritePoint(write.NewPoint(
		taskID,
		map[string]string{},
		map[string]interface{}{
			"productID":         event.ProductID,
			"proportionPercent": event.Command.ProportionPercent,
			"action":            event.Command.Action.String(),
			"task":              event.Task.String(),
		},
		event.CreatedAt,
	))

}

func (d *OrderEventDispatcher) bucketExists(bucket string) bool {
	bucketApi := d.client.BucketsAPI()
	_, err := bucketApi.FindBucketByName(context.Background(), bucket)
	return err == nil
}

func (d *OrderEventDispatcher) Close() error {
	d.writer.Flush()
	d.client.Close()
	return nil
}

// Close closes the dispatcher.
type AnnotationDispatcher struct {
	client influxdb2.Client
	writer api.WriteAPI
}

func NewAnnotationDispatcher(o *Opts) (*AnnotationDispatcher, error) {

	if o.URL == "" {
		return nil, fmt.Errorf("create influx db client: Required field Url is blank")
	}

	if o.Token == "" {
		return nil, fmt.Errorf("create influx db client: Required field Token is blank")
	}

	if o.Org == "" {
		return nil, fmt.Errorf("create influx db client: Required field Url is blank")
	}

	if o.BucketName == "" {
		return nil, fmt.Errorf("create influx db client: Required field TradeBucketName is blank")
	}

	client := influxdb2.NewClient(o.URL, o.Token)

	instance := &AnnotationDispatcher{
		client: client,
		writer: client.WriteAPI(o.Org, o.BucketName),
	}

	if !instance.bucketExists(o.BucketName) {
		return nil, ErrBucketNotExist
	}

	return instance, nil
}

// Dispatch dispatches the given data.
func (d *AnnotationDispatcher) Dispatch(taskID string, data any, createdAt time.Time) {
	influxDataField, _ := mapper.StructToPoint(data)
	d.writer.WritePoint(write.NewPoint(
		taskID,
		map[string]string{},
		influxDataField,
		createdAt,
	))
}

// Close closes the dispatcher.
func (d *AnnotationDispatcher) Close() error {
	d.writer.Flush()
	d.client.Close()
	return nil
}

func (d *AnnotationDispatcher) bucketExists(bucket string) bool {
	bucketApi := d.client.BucketsAPI()
	_, err := bucketApi.FindBucketByName(context.Background(), bucket)
	return err == nil
}
