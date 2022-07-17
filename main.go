package main

import (
	"github.com/tryfix/log"
	"os"
)

func main() {
	logger := log.Constructor.Log(
		log.WithColors(true),
		log.WithLevel("DEBUG"),
		log.WithFilePath(true),
	)
	args := os.Args

	ssiAgent := initAgent(args[1], logger)
	s := newStore(args[1], logger)
	s.init(ssiAgent.ctx)
	defer s.Close()
	initHttpClient(args[1], ssiAgent, s, logger)
}
