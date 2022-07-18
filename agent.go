package main

import (
	"fmt"
	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/issuecredential"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	issuecredential2 "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/issuecredential"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/ws"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	connPkg "github.com/hyperledger/aries-framework-go/pkg/store/connection"
	"github.com/tryfix/log"
	"strconv"
)

type connection struct {
	myDid           string
	theirDid        string
	serviceEndpoint string
}

type agent struct {
	port   int
	client *didexchange.Client
	issuer *issuecredential.Client
	conn   connection
	logger log.Logger
}

func newAgent(port int, logger log.Logger) *agent {
	// inbound transport for agent
	address := fmt.Sprintf("localhost:%d", port+1)
	inbound, err := ws.NewInbound(address, "ws://"+address, "", "")
	if err != nil {
		logger.Fatal(err)
	}

	// aries framework and context initialization
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

	// did-exchange client
	client, err := didexchange.New(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	didActions := make(chan service.DIDCommAction, 1)
	err = client.RegisterActionEvent(didActions)
	if err != nil {
		logger.Fatal(err)
	}
	go service.AutoExecuteActionEvent(didActions)

	// credential issuer
	issuer, err := issuecredential.New(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	a := &agent{port: port + 1, client: client, issuer: issuer, logger: logger}

	issueActions := make(chan service.DIDCommAction)
	err = issuer.RegisterActionEvent(issueActions)
	if err != nil {
		logger.Fatal(err)
	}
	go a.listen(issueActions)

	return a
}

func (a *agent) createInv() (*didexchange.Invitation, error) {
	inv, err := a.client.CreateInvitation(fmt.Sprintf("agent %d", a.port))
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (a *agent) connect(inv *didexchange.Invitation) (*didexchange.Connection, error) {
	connID, err := a.client.HandleInvitation(inv)
	if err != nil {
		return nil, err
	}

	conn, err := a.client.GetConnection(connID)
	if err != nil {
		return nil, err
	}
	a.setConn(conn)

	return conn, nil
}

func (a *agent) getConn() ([]*didexchange.Connection, error) {
	conns, err := a.client.QueryConnections(&didexchange.QueryConnectionsParams{})
	if err != nil {
		return nil, err
	}

	if len(conns) > 0 {
		a.setConn(conns[0])
		return conns, nil
	}

	a.logger.Debug("no connections found")
	return nil, nil
}

// should be called in a more appropriate place in both agents (not in getConn)
func (a *agent) setConn(conn *didexchange.Connection) {
	a.conn.myDid = conn.MyDID
	a.conn.theirDid = conn.TheirDID
	a.conn.serviceEndpoint = conn.ServiceEndPoint
}

func (a *agent) sendOffer() {
	if _, err := a.issuer.SendOffer(&issuecredential.OfferCredential{}, &connPkg.Record{
		MyDID:    a.conn.myDid,
		TheirDID: a.conn.theirDid,
	}); err != nil {
		a.logger.Error(err)
		return
	}
	a.logger.Debug("offer sent to the holder")
}

func (a *agent) listen(issueActions chan service.DIDCommAction) {
	for {
		select {
		case e := <-issueActions:
			piid := e.Properties.All()["piid"].(string)
			if e.Message.Type() == issuecredential2.OfferCredentialMsgTypeV2 {
				if err := a.issuer.AcceptOffer(piid, &issuecredential.RequestCredential{}); err != nil {
					a.logger.Error(err)
					continue
				}
				a.logger.Debug("offer for credential is received")
			}

			if e.Message.Type() == issuecredential2.RequestCredentialMsgTypeV2 {
				if err := a.issuer.AcceptRequest(piid, &issuecredential.IssueCredential{}); err != nil {
					a.logger.Error(err)
					continue
				}
				a.logger.Debug("request for credential is received")
			}

			if e.Message.Type() == issuecredential2.IssueCredentialMsgTypeV2 {
				if err := a.issuer.AcceptCredential(piid, issuecredential.AcceptByFriendlyNames("agent "+strconv.Itoa(a.port))); err != nil {
					a.logger.Error(err)
					continue
				}
				a.logger.Debug("credential has been issued")
			}

			if e.Message.Type() == issuecredential2.ProblemReportMsgTypeV3 {
				if err := a.issuer.AcceptProblemReport(piid); err != nil {
					a.logger.Error(err)
					continue
				}
				a.logger.Error("problem report occurred", e.Message)
			}
		}
	}
}
