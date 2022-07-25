package main

import (
	"github.com/YasiruR/agent/agent"
	"github.com/YasiruR/agent/transport"
	"github.com/tryfix/log"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		log.Fatal(`incorrect number of arguments [./agent <port> <url>]`)
	}

	port, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal(err)
	}

	logger := log.Constructor.Log(log.WithColors(true), log.WithLevel("DEBUG"), log.WithFilePath(true))
	a := agent.New(port, args[2], logger)
	transport.New(port, a, logger).Serve()
}
