package main

import (
	"github.com/tryfix/log"
	"os"
	"strconv"
)

func main() {
	logger := log.Constructor.Log(
		log.WithColors(true),
		log.WithLevel("DEBUG"),
		log.WithFilePath(true),
	)
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		logger.Fatal(err)
	}

	a := newAgent(port, logger)
	newHttpClient(port, a, logger)
}
