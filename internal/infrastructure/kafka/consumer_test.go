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
	}, newSubscribeListenerImpl(ch))
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownConsumer(c *kafka.Consumer[*model_latest.Event]) {
	c.Close()
}


func Test_Consumer(t *testing.T) {

	c := SetupConsumer(make(chan *model_latest.Event, 10))
	defer TeardownConsumer(c)

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	
		err := c.Ping(ctx)
		assert.NoError(t, err)
	})
}