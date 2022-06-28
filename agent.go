package main

import (
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

type agent struct {
	oob struct {
		client *outofbandv2.Client
	}
	vc struct {
		issuer *issuecredential.Client
		stream chan service.DIDCommAction
	}
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

	issuer, issuerActions := initIssuer(ctx, logger)

	a := agent{logger: logger}
	a.vc.issuer = issuer
	a.vc.stream = issuerActions

	//go a.handleVCIssuance()

	// todo test
	oob, err := outofbandv2.New(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	a.oob.client = oob

	//oobDids, err := didexchange.New(ctx)
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//oobEvents := make(chan service.DIDCommAction, 5)
	//
	//err = oobDids.RegisterActionEvent(oobEvents)
	//if err != nil {
	//	logger.Fatal(err)
	//}

	return &a
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
