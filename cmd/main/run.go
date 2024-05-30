package main

import (
	"github.com/Goboolean/core-system.worker/cmd/compose"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Print(compose.MockRun())

}
