package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/aries-framework-go-ext/component/vdr/indy"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/tryfix/log"
	"os"
)

type agent struct {
	client *didexchange.Client
	vdr    *indy.VDR
	logger log.Logger
}

func initAgent(port string, logger log.Logger) *agent {
	framework, err := aries.New()
	if err != nil {
		logger.Fatal(err)
	}

	ctx, err := framework.Context()
	if err != nil {
		logger.Fatal(err)
	}

	client, err := didexchange.New(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	clientActions := make(chan service.DIDCommAction, 1)
	err = client.RegisterActionEvent(clientActions)
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

	logger.Info("agent initialized")

	//go service.AutoExecuteActionEvent(clientActions)
	return &agent{client: client, vdr: vdr, logger: logger}
}

func (a *agent) createInvitation() (didexchange.Invitation, error) {
	inv, err := a.client.CreateInvitation("ssi agent invites a user")
	if err != nil {
		return didexchange.Invitation{}, err
	}

	return *inv, nil
}

func (a *agent) handleInvitation(inv *didexchange.Invitation) (connID string, err error) {
	connID, err = a.client.HandleInvitation(inv)
	if err != nil {
		return
	}

	// create did doc
	// create did for connection
	// create connection request

	return
}

func (a *agent) connection(id string) (didexchange.Connection, error) {
	conn, err := a.client.GetConnection(id)
	if err != nil {
		return didexchange.Connection{}, errors.New(fmt.Sprintf("%s [%s]", err.Error(), id))
	}

	return *conn, nil
}
