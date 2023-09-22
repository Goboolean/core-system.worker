package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)



type Serializer interface {
	Serialize(topic string, v interface{}) ([]byte, error)
}

type ProtoSerializer struct {}
func (s *ProtoSerializer) Serialize(topic string, v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func defaultSerializer() Serializer {
	return &ProtoSerializer{}
}



type Producer struct {
	producer *kafka.Producer
	serial   Serializer
	registry schemaregistry.Client

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// example:
// p, err := NewProducer(&resolver.ConfigMap{
//   "BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
//   "REGISTRY_HOST":  os.Getenv("KAFKA_REGISTRY_HOST"), // optional
// })
func NewProducer(c *resolver.ConfigMap) (*Producer, error) {

	bootstrap_host, err := c.GetStringKey("BOOTSTRAP_HOST")
	if err != nil {
		return nil, err
	}

	registry_host, exists, err := c.GetStringKeyOptional("REGISTRY_HOST")
	if err != nil {
		return nil, err
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   bootstrap_host,
		"acks":                -1,
		"go.delivery.reports": true,
	})

	ctx, cancel := context.WithCancel(context.Background())

	instance := &Producer{
		producer: p,
		wg: sync.WaitGroup{},
		ctx: ctx,
		cancel: cancel,
	}

	if exists {
		r, err := schemaregistry.NewClient(schemaregistry.NewConfig(registry_host))
		if err != nil {
			return nil, err
		}		

		s, err := protobuf.NewSerializer(r, serde.ValueSerde, protobuf.NewSerializerConfig())
		if err != nil {
			return nil, err
		}
		instance.serial = s

	} else {
		instance.serial = defaultSerializer()
	}

	instance.traceEvent(ctx, &instance.wg)
	return instance, nil
}


// The schema argument can be provided by protobuf generated struct,
// initialized with default values.
// This function returns the ID of the schema in the registry.
// TODO: implement a returning ID logic.
func (p *Producer) Register(topic string, schema proto.Message) (int64, error) {
	_, err := p.serial.Serialize(topic, schema)
	return 0, err
}


func (p *Producer) Produce(topic string, msg proto.Message) error {
	payload, err := p.serial.Serialize(topic, msg)
	if err != nil {
		return err
	}

	if err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Value:          payload,
	}, nil); err != nil {
		return err
	}

	return nil
}


func (p *Producer) Flush(ctx context.Context) (int, error) {

	deadline, ok := ctx.Deadline()
	if !ok {
		return 0, ErrDeadlineSettingRequired
	}

	left := p.producer.Flush(int(time.Until(deadline).Milliseconds()))
	if left != 0 {
		return left, ErrFailedToFlush
	}

	return 0, nil
}


func (p *Producer) traceEvent(ctx context.Context, wg *sync.WaitGroup) {

	go func() {
		wg.Add(1)
		defer wg.Done()

		for e := range p.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.WithFields(log.Fields{
						"topic": *ev.TopicPartition.Topic,
						"data":  ev.Value,
						"error":  ev.TopicPartition.Error,
					}).Error("Producer failed to deliver event to kafka")
				} else {
					log.WithFields(log.Fields{
						"topic": *ev.TopicPartition.Topic,
						"data":  ev.Value,
						"partition":  ev.TopicPartition.Partition,
						"offset": ev.TopicPartition.Offset,
					}).Trace("Producer delivered event to kafka")
				}
			}
		}
	}()
}


func (p *Producer) Close() {
	p.producer.Close()
	p.cancel()
	p.wg.Wait()
}


func (p *Producer) Ping(ctx context.Context) error {
	// It requires ctx to be deadline set, otherwise it will return error
	// It will return error if there is no response within deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		return ErrDeadlineSettingRequired
	}

	remaining := time.Until(deadline)
	_, err := p.producer.GetMetadata(nil, true, int(remaining.Milliseconds()))
	if err != nil {
		return err
	}

	if p.registry != nil {
		_, err := p.registry.GetAllVersions("schema")
		if err != nil {
			return err
		}
	}

	return nil
}