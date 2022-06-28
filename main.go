package main

import (
	"github.com/tryfix/log"
	"os"
)

func main() {
	logger := log.Constructor.Log(log.WithFilePath(true))
	args := os.Args
	initHttpClient(args[1], initAgent(logger), logger)
}
