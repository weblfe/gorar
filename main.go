package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/gorar/commands"
)

var version = "0.0.1"

func main() {
	if err := commands.New(version).Execute(); err != nil {
		log.WithField("error", err).Errorln()
	}
}
