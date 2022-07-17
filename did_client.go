package main

import (
	"fmt"
	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/ws"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/tryfix/log"
)

type agent struct {
	port int
	didexchange.Client
	logger log.Logger
}

func newAgent(port int, logger log.Logger) *agent {
	address := fmt.Sprintf("localhost:%d", port+1)
	inbound, err := ws.NewInbound(address, "ws://"+address, "", "")
	if err != nil {
		logger.Fatal(err)
	}

	fw, err := aries.New(
		aries.WithInboundTransport(inbound),
		aries.WithOutboundTransports(ws.NewOutbound()),
		aries.WithStoreProvider(mem.NewProvider()),
		aries.WithProtocolStateStoreProvider(mem.NewProvider()),
	)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, err := fw.Context()
	if err != nil {
		logger.Fatal(err)
	}

	client, err := didexchange.New(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	actions := make(chan service.DIDCommAction, 1)
	err = client.RegisterActionEvent(actions)
	if err != nil {
		logger.Fatal(err)
	}

	go service.AutoExecuteActionEvent(actions)
	return &agent{port: port + 1, Client: *client, logger: logger}
}

func (a *agent) createInv() (*didexchange.Invitation, error) {
	inv, err := a.Client.CreateInvitation(fmt.Sprintf("agent %d", a.port))
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (a *agent) connect(inv *didexchange.Invitation) (*didexchange.Connection, error) {
	connID, err := a.Client.HandleInvitation(inv)
	if err != nil {
		return nil, err
	}

	conn, err := a.Client.GetConnection(connID)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (a *agent) getConn() ([]*didexchange.Connection, error) {
	conns, err := a.Client.QueryConnections(&didexchange.QueryConnectionsParams{})
	if err != nil {
		return nil, err
	}
	return conns, nil
}

//type inboundTransport struct {
//	port int
//}
//
//func newInboundTransport(port int) *inboundTransport {
//	return &inboundTransport{port: port + 1}
//}
//
//func (i *inboundTransport) Start(prov transport.Provider) error {
//	return nil
//}
//
//func (i *inboundTransport) Stop() error {
//	return nil
//}
//
//func (i *inboundTransport) Endpoint() string {
//	return fmt.Sprintf("http://localhost:%d", i.port)
//}
