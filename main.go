package main

import (
	"flag"
	"github.com/YasiruR/agent/agent"
	"github.com/YasiruR/agent/transport"
	"github.com/tryfix/log"
)

func main() {
	port, url := parseArgs()
	logger := log.Constructor.Log(log.WithColors(true), log.WithLevel("DEBUG"), log.WithFilePath(true))
	a := agent.New(port, url, logger)
	transport.New(port, a, logger).Serve()
}

func parseArgs() (port int, url string) {
	p := flag.Int(`port`, 0, `port of the controller`)
	u := flag.String(`url`, ``, `url of the agent`)
	flag.Parse()

	if *p == 0 {
		log.Fatal(`port for controller must be specified`)
	}

	return *p, *u
}
