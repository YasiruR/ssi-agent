package main

import (
	"github.com/hyperledger/aries-framework-go-ext/component/vdr/indy"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"
	"github.com/tryfix/log"
	"os"
)

type agent struct {
	ctx       *context.Provider
	oobClient *outofband.Client
	vdr       *indy.VDR
	logger    log.Logger
}

func initAgent(port string, logger log.Logger) *agent {
	// initialize aries instance and context
	framework, err := aries.New()
	if err != nil {
		logger.Fatal(err)
	}

	ctx, err := framework.Context()
	if err != nil {
		logger.Fatal(err)
	}

	// initiate indy vdr
	pwd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}

	vdr, err := indy.New(`sov`, indy.WithIndyVDRGenesisFile(pwd+`/src/genesis.json`))
	if err != nil {
		logger.Fatal(err)
	}

	// initialize out-of-band client
	oobClient, err := outofband.New(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("agent initialized")
	return &agent{ctx: ctx, oobClient: oobClient, vdr: vdr, logger: logger}
}

func (a *agent) createInv() (*outofband.Invitation, error) {
	inv, err := a.oobClient.CreateInvitation(nil)
	if err != nil {
		return nil, err
	}

	return inv, nil
}

func (a *agent) acceptInv(inv *outofband.Invitation) (string, error) {
	connID, err := a.oobClient.AcceptInvitation(inv, "agent accepts invitation")
	if err != nil {
		return "", err
	}

	a.logger.Debug("out-of-band invitation accepted via agent", connID)
	return connID, nil
}
