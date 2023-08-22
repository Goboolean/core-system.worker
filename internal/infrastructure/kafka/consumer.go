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
	"google.golang.org/protobuf/reflect/protoreflect"
)



type Deserializer interface {
	DeserializeInto(topic string, payload []byte, msg interface{}) error
}

type ProtoDeserializer struct {}

func (s *ProtoDeserializer) DeserializeInto(topic string, payload []byte, msg interface{}) error {
	return proto.Unmarshal(payload, msg.(proto.Message))
}

func newDeserializer() Deserializer {
	return &ProtoDeserializer{}
}


type Consumer struct {
	consumer *kafka.Consumer
	deserial Deserializer
}

func NewConsumer(c *resolver.ConfigMap) (*Consumer, error) {

	bootstrap_host, err := c.GetStringKey("BOOTSTRAP_HOST")
	if err != nil {
		return nil, err
	}

	group_id, err := c.GetStringKey("GROUP_ID")
	if err != nil {
		return nil, err
	}

	con, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":   bootstrap_host,
		"group.id":            group_id,
	})

	instance := &Consumer{
		consumer: con,
	}

	registry_url, exists, err := c.GetStringKeyOptional("REGISTRY_URL")
	if err != nil {
		return nil, err
	}

	if exists {
		sr, err := schemaregistry.NewClient(schemaregistry.NewConfig(registry_url))
		if err != nil {
			return nil, err
		}

		d, err := protobuf.NewDeserializer(sr, serde.ValueSerde, protobuf.NewDeserializerConfig())
		if err != nil {
			return nil, err
		}

		instance.deserial = d
	} else {
		instance.deserial = newDeserializer()
	}

	return instance, nil
}


func (c *Consumer) Subscribe(topic string, schema protoreflect.MessageType) error {

	_, ok := c.deserial.(*protobuf.Deserializer)
	if ok {
		if err := c.deserial.(*protobuf.Deserializer).ProtoRegistry.RegisterMessage(schema); err != nil {
			return err
		}
	}

	if err := c.consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		return err
	}
	return nil
}


func (c *Consumer) Close() {
	c.consumer.Close()
}


func (c *Consumer) Ping(ctx context.Context) error {
	// It requires ctx to be deadline set, otherwise it will return error
	// It will return error if there is no response within deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		return ErrDeadlineSettingRequired
	}

	remaining := time.Until(deadline)
	_, err := c.consumer.GetMetadata(nil, true, int(remaining.Milliseconds()))
	return err
}