package kafka

import (
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/proto"
)



type Deserializer interface {
	Deserialize(topic string, payload []byte) (interface{}, error)
}

type ProtoDeserializer struct {}

func (s *ProtoDeserializer) Deserialize(topic string, payload []byte) (interface{}, error) {
	var v proto.Message
	return v, proto.Unmarshal(payload, v)
}

func newDeserializer() Deserializer {
	return &ProtoDeserializer{}
}


type Consumer struct {
	consumer *kafka.Consumer
	deserial *protobuf.Deserializer
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
	}

	return instance, nil
}


func (c *Consumer) Subscribe(topic string, schema proto.Message) error {
	if err := c.deserial.ProtoRegistry.RegisterMessage(schema.ProtoReflect().Type()); err != nil {
		return err
	}

	return nil
}