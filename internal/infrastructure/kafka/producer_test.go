package kafka_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	model_latest "github.com/Goboolean/core-system.worker/api/kafka/model.latest"
	"github.com/Goboolean/core-system.worker/internal/infrastructure/kafka"
	"github.com/stretchr/testify/assert"
)





func SetupProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
	})
	if err != nil {
		panic(err)
	}
	return p
}

func SetupProducerWithRegistry() *kafka.Producer {
	p, err := kafka.NewProducer(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
		"REGISTRY_HOST":  os.Getenv("KAFKA_REGISTRY_HOST"),
	})
	if err != nil {
		panic(err)
	}
	return p
}

func TeardownProducer(p *kafka.Producer) {
	mutex.Lock()
	defer mutex.Unlock()
	p.Close()
}



func Test_Producer(t *testing.T) {

	p := SetupProducer()
	defer TeardownProducer(p)

	const topic = "default-topic"
	var event = &model_latest.Event{
		EventUuid: "test-uuid",
	}

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	
		err := p.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Register", func(t *testing.T) {	
		_, err := p.Register(topic, &model_latest.Event{})
		assert.NoError(t, err)
	})

	t.Run("ProduceMessage", func(t *testing.T) {
		err := p.Produce(topic, event)
		assert.NoError(t, err)
	})

	t.Run("Flush", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		count, err := p.Flush(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}


func Test_ProducerWithRegistry(t *testing.T) {
	t.Skip("Skip this test because of the registry is not ready.")

	p := SetupProducerWithRegistry()
	defer TeardownProducer(p)

	const topic = "default-topic"
	var event = &model_latest.Event{
		EventUuid: "test-uuid",
	}

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	
		err := p.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Register", func(t *testing.T) {	
		_, err := p.Register(topic, &model_latest.Event{})
		assert.NoError(t, err)
	})

	t.Run("ProduceMessage", func(t *testing.T) {
		err := p.Produce(topic, event)
		assert.NoError(t, err)
	})

	t.Run("Flush", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		count, err := p.Flush(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}