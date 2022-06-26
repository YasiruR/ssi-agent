package main

import (
	"github.com/tryfix/log"
	"os"
)

func main() {
	logger := log.NewLog().Log()
	args := os.Args
	initHttpClient(args[1], initAgent(logger), logger)
}
