package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// !!! Import this package at main.
// _ "github.com/Goboolean/fetch-server.v1/internal/util/logger"
func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
