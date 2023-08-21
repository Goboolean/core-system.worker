package kafka

import (
	"context"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/proto"
)



type Serializer interface {
	Serialize(topic string, v interface{}) ([]byte, error)
}

type ProtoSerializer struct {}
func (s *ProtoSerializer) Serialize(topic string, v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}



type Producer struct {
	producer *kafka.Producer
	serial   Serializer
	registry schemaregistry.Client
}

func NewProducer(c *resolver.ConfigMap) (*Producer, error) {

	bootstrap_host, err := c.GetStringKey("BOOTSTRAP_HOST")
	if err != nil {
		return nil, err
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   bootstrap_host,
		"acks":                -1,
		"go.delivery.reports": true,
	})

	instance := &Producer{
		producer: p,
	}

	registry_host, exists, err := c.GetStringKeyOptional("REGISTRY_HOST")
	if err != nil {
		return nil, err
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
		instance.serial = &ProtoSerializer{}
	}

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


func (p *Producer) Close() {
	p.producer.Close()
}