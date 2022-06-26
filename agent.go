package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/tryfix/log"
)

type agent struct {
	client *didexchange.Client
	logger log.Logger
}

func initAgent(logger log.Logger) *agent {
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

	//go service.AutoExecuteActionEvent(clientActions)

	return &agent{client: client, logger: logger}
}

func (a *agent) createInvitation() (didexchange.Invitation, error) {
	inv, err := a.client.CreateInvitation("ssi agent invites a user")
	if err != nil {
		return didexchange.Invitation{}, err
	}

	connID, err := a.handleInvitation(inv)
	if err != nil {
		return didexchange.Invitation{}, err
	}

	conn, err := a.connection(connID)
	if err != nil {
		return didexchange.Invitation{}, err
	}
	fmt.Println("(sender) connection details: ", conn.State, conn.ConnectionID, conn.InvitationID)

	return *inv, nil
}

func (a *agent) handleInvitation(inv *didexchange.Invitation) (connID string, err error) {
	connID, err = a.client.HandleInvitation(inv)
	if err != nil {
		return
	}

	return
}

func (a *agent) connection(id string) (didexchange.Connection, error) {
	conn, err := a.client.GetConnection(id)
	if err != nil {
		return didexchange.Connection{}, errors.New(fmt.Sprintf("%s [%s]", err.Error(), id))
	}

	return *conn, nil
}
