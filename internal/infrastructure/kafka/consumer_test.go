package kafka_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/api/kafka/model.latest"
	"github.com/Goboolean/core-system.worker/internal/infrastructure/kafka"
	"github.com/stretchr/testify/assert"
)


type SubscribeListenerImpl struct {
	channel chan *model_latest.Event
}

func newSubscribeListenerImpl(ch chan *model_latest.Event) *SubscribeListenerImpl {
	return &SubscribeListenerImpl{
		channel: ch,
	}
}

func (s *SubscribeListenerImpl) OnReceiveMessage(ctx context.Context, msg *model_latest.Event) error {
	s.channel <- msg
	return nil
}



func SetupConsumer(ch chan *model_latest.Event) *kafka.Consumer[*model_latest.Event] {

	c, err := kafka.NewConsumer[*model_latest.Event](&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
		"GROUP_ID": "TEST_GROUP",
		"PROCESSOR_COUNT": 2,
	}, newSubscribeListenerImpl(ch))
	if err != nil {
		panic(err)
	}
	return c
}

func SetupConsumerWithRegistry(ch chan *model_latest.Event) *kafka.Consumer[*model_latest.Event] {

	c, err := kafka.NewConsumer[*model_latest.Event](&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
		"REGISTRY_HOST": os.Getenv("KAFKA_REGISTRY_HOST"),
		"GROUP_ID": "TEST_GROUP",
		"PROCESSOR_COUNT": 2,
	}, newSubscribeListenerImpl(ch))
	if err != nil {
		panic(err)
	}
	return c
}


func TeardownConsumer(c *kafka.Consumer[*model_latest.Event]) {
	mutex.Lock()
	defer mutex.Unlock()
	c.Close()
}


func Test_Consumer(t *testing.T) {

	const topic = "default-topic"
	var event = &model_latest.Event{}

	var channel = make(chan *model_latest.Event, 10)

	c := SetupConsumer(channel)
	defer TeardownConsumer(c)
	p := SetupProducer()
	defer TeardownProducer(p)

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	
		err := c.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Subscribe", func(t *testing.T) {
		err := c.Subscribe(topic, event.ProtoReflect().Type())
		assert.NoError(t, err)
	})

	t.Run("ConsumeMessage", func(t *testing.T) {
		err := p.Produce(topic, event)
		assert.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		select {
		case <-ctx.Done():
			assert.Fail(t, "timeout")
			break
		case <-channel:
			break
		}
	})
}