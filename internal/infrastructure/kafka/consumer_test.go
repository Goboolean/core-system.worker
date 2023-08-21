package kafka_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/infrastructure/kafka"
	"github.com/stretchr/testify/assert"
)





func SetupConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
		"GROUP_ID": "TEST_GROUP",
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownConsumer(c *kafka.Consumer) {
	c.Close()
}


func Test_Consumer(t *testing.T) {

	c := SetupConsumer()
	defer TeardownConsumer(c)

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	
		err := c.Ping(ctx)
		assert.NoError(t, err)
	})
}