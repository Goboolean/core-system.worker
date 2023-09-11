package kafka_test

import (
	"os"
	"sync"
	"testing"

	_ "github.com/Goboolean/core-system.worker/internal/util/env"
	log "github.com/sirupsen/logrus"
)

var mutex = &sync.Mutex{}



func TestMain(m *testing.M) {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetLevel(log.TraceLevel)

	code := m.Run()
	os.Exit(code)
}