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

func SetupConsumerWithGroup(ch chan *model_latest.Event, group string) *kafka.Consumer[*model_latest.Event] {

	c, err := kafka.NewConsumer[*model_latest.Event](&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
		"GROUP_ID": group,
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

	const topic = "test-consumer"
	var event = &model_latest.Event{
		EventUuid: "test-uuid",
	}

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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 7)
		defer cancel()

		select {
		case <-ctx.Done():
			assert.Fail(t, "failed to receive all message before timeout")
			return
		case <-channel:
			break
		}
	})
}


func Test_ConsumerWithRegistry(t *testing.T) {
	t.Skip("Skip this test because of the registry is not ready.")

	const topic = "test-consumer-with-registry"
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
		_, err := p.Register(topic, event)
		assert.NoError(t, err)

		err = p.Produce(topic, event)
		assert.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		select {
		case <-ctx.Done():
			assert.Fail(t, "failed to receive all message before timeout")
			return
		case <-channel:
			break
		}
	})
}


func Test_ConsumeSameGroup(t *testing.T) {

	const topic = "test-consumer-same-group"
	const count = 3
	var event = &model_latest.Event{EventUuid: "test-uuid"}

	var ch1 = make(chan *model_latest.Event, 10)
	var ch2 = make(chan *model_latest.Event, 10)

	c1 := SetupConsumerWithGroup(ch1, "TEST_GROUP")
	defer TeardownConsumer(c1)
	c2 := SetupConsumerWithGroup(ch2, "TEST_GROUP")
	defer TeardownConsumer(c2)

	p := SetupProducer()
	defer TeardownProducer(p)

	t.Run("Subscribe", func(t *testing.T) {
		err := c1.Subscribe(topic, event.ProtoReflect().Type())
		assert.NoError(t, err)

		err = c2.Subscribe(topic, event.ProtoReflect().Type())
		assert.NoError(t, err)
	})

	t.Run("ConsumeMessage", func(t *testing.T) {
		for i := 0; i < count; i++ {
			err := p.Produce(topic, event)
			assert.NoError(t, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		<- ctx.Done()
		assert.Equal(t, count, len(ch1) + len(ch2))
	})
}




func Test_ConsumeDifferentGroup(t *testing.T) {

	const topic = "test-consumer-different-group"
	const count = 3
	var event = &model_latest.Event{EventUuid: "test-uuid"}

	var ch_a = make(chan *model_latest.Event, 10)
	var ch_b = make(chan *model_latest.Event, 10)

	c_a := SetupConsumerWithGroup(ch_a, "TEST_GROUP_A")
	defer TeardownConsumer(c_a)
	c_b := SetupConsumerWithGroup(ch_b, "TEST_GROUP_B")
	defer TeardownConsumer(c_b)

	p := SetupProducer()
	defer TeardownProducer(p)

	t.Run("Subscribe", func(t *testing.T) {
		err := c_a.Subscribe(topic, event.ProtoReflect().Type())
		assert.NoError(t, err)

		err = c_b.Subscribe(topic, event.ProtoReflect().Type())
		assert.NoError(t, err)
	})

	t.Run("ConsumeMessage", func(t *testing.T) {
		for i := 0; i < count; i++ {
			err := p.Produce(topic, event)
			assert.NoError(t, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		<- ctx.Done()
		assert.Equal(t, count, len(ch_a))
		assert.Equal(t, count, len(ch_b))
	})
}
