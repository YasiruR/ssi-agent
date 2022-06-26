package main

import (
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"log"
)

type agent struct {
	client *didexchange.Client
}

func initAgent() *agent {
	framework, err := aries.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx, err := framework.Context()
	if err != nil {
		log.Fatal(err)
	}

	client, err := didexchange.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	clientActions := make(chan service.DIDCommAction, 1)
	err = client.RegisterActionEvent(clientActions)
	if err != nil {
		log.Fatal(err)
	}

	go service.AutoExecuteActionEvent(clientActions)

	return &agent{client: client}
}

func (a *agent) CreateInvitation() (*didexchange.Invitation, error) {
	inv, err := a.client.CreateInvitation("ssi agent invites a user")
	if err != nil {
		return nil, err
	}

	return inv, nil
}

func (a *agent) HandleInvitation(inv *didexchange.Invitation) (connID string, err error) {
	return a.client.HandleInvitation(inv)
}
