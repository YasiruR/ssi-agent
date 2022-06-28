package main

import (
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/issuecredential"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofbandv2"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	issuecredential2 "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/issuecredential"
	outofbandv22 "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/outofbandv2"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"
	connection2 "github.com/hyperledger/aries-framework-go/pkg/store/connection"
	"github.com/tryfix/log"
)

type connection struct {
	id          string
	ownDID      string
	receiverDID string
}

type agent struct {
	did struct {
		client *didexchange.Client
		stream chan service.DIDCommAction
	}
	oob struct {
		client *outofbandv2.Client
	}
	vc struct {
		issuer *issuecredential.Client
		stream chan service.DIDCommAction
	}
	connections map[string]connection
	logger      log.Logger
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

	client, clientActions := initDidClient(ctx, logger)
	issuer, issuerActions := initIssuer(ctx, logger)

	a := agent{logger: logger, connections: map[string]connection{}}
	a.did.client = client
	a.did.stream = clientActions
	a.vc.issuer = issuer
	a.vc.stream = issuerActions

	//go a.handleVCIssuance()

	// todo test
	oob, err := outofbandv2.New(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	a.oob.client = oob

	return &a
}

func initDidClient(ctx *context.Provider, logger log.Logger) (*didexchange.Client, chan service.DIDCommAction) {
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

	return client, clientActions
}

func initIssuer(ctx *context.Provider, logger log.Logger) (*issuecredential.Client, chan service.DIDCommAction) {
	issuer, err := issuecredential.New(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	issuerActions := make(chan service.DIDCommAction)
	err = issuer.RegisterActionEvent(issuerActions)
	if err != nil {
		logger.Fatal(err)
	}

	return issuer, issuerActions
}

func (a *agent) createInvitation() (*outofbandv22.Invitation, error) {
	inv, err := a.oob.client.CreateInvitation()
	if err != nil {
		return nil, err
	}

	return inv, nil
}

func (a *agent) acceptInvitation(inv *outofbandv22.Invitation) error {
	connID, err := a.oob.client.AcceptInvitation(inv)
	if err != nil {
		return err
	}

	a.logger.Debug("accepted invitation conn id: ", connID)
	return nil
}

//func (a *agent) createInvitation() (*didexchange.Invitation, error) {
//	inv, err := a.did.client.CreateInvitation("ssi agent invites a user")
//	if err != nil {
//		return nil, err
//	}
//
//	connID, err := a.handleInvitation(inv)
//	if err != nil {
//		return nil, err
//	}
//
//	conn, err := a.connection(connID)
//	if err != nil {
//		return nil, err
//	}
//	a.logger.Debug("sender connection details (create invitation): ", conn.State, conn.ConnectionID, conn.MyDID, conn.TheirDID)
//	return inv, nil
//}

//func (a *agent) handleInvitation(inv *didexchange.Invitation) (connID string, err error) {
//	connID, err = a.did.client.HandleInvitation(inv)
//	if err != nil {
//		return
//	}
//
//	//conn, err := a.connection(connID)
//	//if err != nil {
//	//	a.logger.Error(err)
//	//}
//	//a.logger.Debug("receiver connection details (handle invitation): ", conn.State, conn.ConnectionID, conn.MyDID, conn.TheirDID)
//
//	return
//}

//func (a *agent) acceptInvitation(inv *didexchange.Invitation) error {
//	connID, err := a.handleInvitation(inv)
//	if err != nil {
//		return err
//	}
//
//	err = a.did.client.AcceptInvitation(connID, inv.RecipientKeys[0], inv.Label)
//	if err != nil {
//		return err
//	}
//
//	//conn, err := a.connection(connID)
//	//if err != nil {
//	//	return err
//	//}
//
//	//a.logger.Debug("accepting invitation")
//	//fmt.Println("conn id: ", conn.ConnectionID)
//	//fmt.Println("inv id: ", conn.InvitationID)
//	//fmt.Println("my did: ", conn.MyDID)
//	//fmt.Println("their did: ", conn.TheirDID)
//	//fmt.Println("rec keys: ", conn.RecipientKeys)
//	//fmt.Println("routing keys: ", conn.RoutingKeys)
//	//fmt.Println("state: ", conn.State)
//
//	return nil
//}

//func (a *agent) connection(connID string) (*didexchange.Connection, error) {
//	conn, err := a.did.client.GetConnection(connID)
//	if err != nil {
//		return &didexchange.Connection{}, errors.New(fmt.Sprintf("%s [%s]", err.Error(), connID))
//	}
//
//	return conn, nil
//}

//func (a *agent) connectionByState(connID string) (*didexchange.Connection, error) {
//	conn, err := a.did.client.GetConnectionAtState(connID, "requested")
//	if err != nil {
//		return &didexchange.Connection{}, errors.New(fmt.Sprintf("%s [%s]", err.Error(), connID))
//	}
//
//	return conn, nil
//}

func (a *agent) sendVCOffer(record *connection2.Record) error {
	out, err := a.vc.issuer.SendOffer(&issuecredential.OfferCredential{}, record)
	if err != nil {
		return err
	}

	a.logger.Debug("issuer sent offer to the holder", out)
	return nil
}

func (a *agent) handleVCIssuance() {
	for {
		select {
		case event := <-a.vc.stream:
			piid, ok := event.Properties.All()["piid"].(string)
			if !ok {
				a.logger.Error("invalid event", event.Properties.All()["piid"])
			}

			a.logger.Debug("message received for VC issuance", event.Message.Type())

			switch event.Message.Type() {
			case issuecredential2.ProposeCredentialMsgTypeV2:
				if err := a.vc.issuer.AcceptProposal(piid, &issuecredential.OfferCredential{}); err != nil {
					a.logger.Error(err)
					continue
				}
			case issuecredential2.OfferCredentialMsgTypeV2:
				if err := a.vc.issuer.AcceptOffer(piid, &issuecredential.RequestCredential{}); err != nil {
					a.logger.Error(err)
					continue
				}
			case issuecredential2.RequestCredentialMsgTypeV2:
				if err := a.vc.issuer.AcceptRequest(piid, &issuecredential.IssueCredential{}); err != nil {
					a.logger.Error(err)
					continue
				}
			case issuecredential2.IssueCredentialMsgTypeV2:
				if err := a.vc.issuer.AcceptCredential(piid); err != nil {
					a.logger.Error(err)
					continue
				}
			case issuecredential2.ProblemReportMsgTypeV2:
				if err := a.vc.issuer.AcceptProblemReport(piid); err != nil {
					a.logger.Error(err)
					continue
				}
			}
		}
	}
}
