package main

import (
	"flag"
	"github.com/YasiruR/agent/agent"
	agentServer "github.com/YasiruR/agent/transport/agent"
	webhookServer "github.com/YasiruR/agent/transport/webhook"
	"github.com/tryfix/log"
)

func main() {
	controllerPort, webhookPort, url := parseArgs()
	logger := log.Constructor.Log(log.WithColors(true), log.WithLevel("DEBUG"), log.WithFilePath(true))

	a := agent.New(controllerPort, url, logger)
	go webhookServer.New(webhookPort, a, logger).Serve()
	agentServer.New(controllerPort, a, logger).Serve()
}

func parseArgs() (controllerPort, webhookPort int, url string) {
	cp := flag.Int(`controller_port`, 0, `port of the controller`)
	wp := flag.Int(`webhook_port`, 0, `port of the webhook processor`)
	u := flag.String(`agent_url`, ``, `url of the agent`)
	flag.Parse()

	if *cp == 0 {
		log.Fatal(`port for controller must be specified`)
	}

	if *wp == 0 {
		log.Fatal(`port for webhook processor must be specified`)
	}

	return *cp, *wp, *u
}
