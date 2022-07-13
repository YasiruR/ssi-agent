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
	initHttpClient(args[1], initAgent(args[1], logger), logger)
}
