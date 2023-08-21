package kafka_test

import (
	"os"
	"testing"
	_ "github.com/Goboolean/core-system.worker/internal/util/env"
)



func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}