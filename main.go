package main

import (
	"flag"
	"fmt"
	"github.com/YasiruR/agent/agent"
	agentServer "github.com/YasiruR/agent/transport/agent"
	webhookServer "github.com/YasiruR/agent/transport/webhook"
	"github.com/tryfix/log"
	"strconv"
)

func main() {
	name, controllerPort, webhookPort, url := parseArgs()
	logger := log.Constructor.Log(log.WithColors(true), log.WithLevel("DEBUG"), log.WithFilePath(true))

	a := agent.New(name, url, logger)
	go webhookServer.New(webhookPort, a, logger).Serve()
	agentServer.New(controllerPort, a, logger).Serve()
}

func parseArgs() (name string, controllerPort, webhookPort int, url string) {
	l := flag.String(`label`, ``, `label of the agent`)
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

	if *l == `` {
		*l = strconv.Itoa(*cp)
		log.Info(fmt.Sprintf(`agent label is set to the controller port [%d] since not provided explicitly`, *cp))
	}

	return *l, *cp, *wp, *u
}
