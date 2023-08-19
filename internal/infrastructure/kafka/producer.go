package kafka

import (
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	log "github.com/sirupsen/logrus"
)




type Producer struct {
	producer *kafka.Producer
	serial   serde.Serializer
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

	registry_url, exists, err := c.GetStringKeyOptional("REGISTRY_URL")
	if err != nil {
		return nil, err
	}

	var (
		r schemaregistry.Client
		s serde.Serializer
	)

	if exists {
		r, err = schemaregistry.NewClient(schemaregistry.NewConfig(registry_url))
		if err != nil {
			return nil, err
		}

		s, err = protobuf.NewSerializer(r, serde.ValueSerde, protobuf.NewSerializerConfig())
		if err != nil {
			return nil, err
		}
	}

	instance := &Producer{
		producer: p,
		serial:   s,
	}

	return instance, nil
}


