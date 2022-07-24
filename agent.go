package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go-ext/component/vdr/indy"
	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/issuecredential"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/decorator"
	issuecredential2 "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/issuecredential"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/ws"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	connPkg "github.com/hyperledger/aries-framework-go/pkg/store/connection"
	"github.com/hyperledger/aries-framework-go/spi/storage"
	"github.com/tryfix/log"
	"os"
	"strconv"
)

// did exchange structs
type connection struct {
	myDid           string
	theirDid        string
	serviceEndpoint string
}

type agent struct {
	port   int
	client *didexchange.Client
	issuer *issuecredential.Client
	vdr    *indy.VDR
	store  storage.Store
	conn   connection
	logger log.Logger
}

// preview credential structs
type attribute struct {
	Name     string `json:"name"`
	MimeType string `json:"mime-type"`
	Value    string `json:"value"`
}

type previewCred struct {
	Type       string      `json:"@type"`
	Attributes []attribute `json:"attributes"`
}

// issue credential structs
type jsonData struct {
	Context           []string          `json:"@context"`
	IssuanceDate      string            `json:"issuanceDate"`
	Issuer            map[string]string `json:"issuer"`
	ReferenceNumber   int64             `json:"referenceNumber"`
	Type              []string          `json:"type"`
	CredentialSubject struct {
		ID string `json:"id"`
	} `json:"credentialSubject"`
}

func newAgent(port int, logger log.Logger) *agent {
	// indy vdr
	pwd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}

	vdr, err := indy.New(`indy`, indy.WithIndyVDRGenesisFile(pwd+`/src/genesis.json`))
	if err != nil {
		logger.Fatal(err)
	}

	// store
	storeProv := mem.NewProvider()
	wallet, err := storeProv.OpenStore(`wallet`)
	if err != nil {
		logger.Fatal(err)
	}

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
		aries.WithStoreProvider(storeProv),
		aries.WithProtocolStateStoreProvider(mem.NewProvider()),
		aries.WithVDR(vdr),
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

	a := &agent{port: port + 1, client: client, issuer: issuer, vdr: vdr, store: wallet, logger: logger}

	issueActions := make(chan service.DIDCommAction)
	err = issuer.RegisterActionEvent(issueActions)
	if err != nil {
		logger.Fatal(err)
	}
	go a.listen(issueActions)

	a.read()

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
	rec := connPkg.Record{MyDID: a.conn.myDid, TheirDID: a.conn.theirDid}
	offer := issuecredential.OfferCredential{
		Type: "https://didcomm.org/issue-credential/2.0/offer-credential",
		ID:   uuid.New().String(),
		CredentialPreview: previewCred{
			Type: "https://didcomm.org/issue-credential/2.0/credential-preview",
			Attributes: []attribute{
				{Name: "first_name", Value: "Alan"},
				{Name: "role", Value: "developer"},
				{Name: "country", Value: "Norway"},
			},
		},
		Comment: "initial credential offer from agent " + strconv.Itoa(a.port),
	}

	if _, err := a.issuer.SendOffer(&offer, &rec); err != nil {
		a.logger.Error(err)
		return
	}
	a.logger.Debug(fmt.Sprintf("offer sent to the holder \n[rec: %v] \n[offer: %v]", rec, offer))
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

			attach := decorator.GenericAttachment{}
			attach.ID = uuid.New().String()
			attach.Data.JSON = jsonData{
				Context:         []string{"https://www.w3.org/2018/credentials/v1", "https://www.w3.org/2018/credentials/examples/v1"},
				IssuanceDate:    "2010-01-01T19:23:24Z",
				Issuer:          map[string]string{"id": uuid.New().String()},
				ReferenceNumber: 83294847,
				Type:            []string{"VerifiableCredential"},
				CredentialSubject: struct {
					ID string `json:"id"`
				}{"initial-verifiable-credential"},
			}

			if e.Message.Type() == issuecredential2.RequestCredentialMsgTypeV2 {
				cred := issuecredential.IssueCredential{
					Type:        "https://didcomm.org/issue-credential/2.0/issue-credential",
					ID:          uuid.New().String(),
					Comment:     "agent " + strconv.Itoa(a.port) + " issuing credential",
					Attachments: []decorator.GenericAttachment{attach},
				}

				if err := a.issuer.AcceptRequest(piid, &cred); err != nil {
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
				a.logger.Debug("credential is received")

				bytes, err := a.store.Get(piid)
				if err != nil {
					a.logger.Error("getting vc from store failed", err)
					continue
				}
				a.logger.Debug("fetched credential from store", string(bytes))
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

// todo remove
func (a *agent) read() {
	status, err := a.vdr.Client.GetPoolStatus()
	if err != nil {
		a.logger.Fatal(err)
	}
	a.logger.Info("status fetched from pool", status)

	//nym, err := a.vdr.Client.GetNym("ByuET4QKgGkaiUYGuvS6x3wADnmbYmtSHBAJngaQRhNL")
	//if err != nil {
	//	a.logger.Error(err)
	//}
	//a.logger.Info("nym fetched", nym)

	//doc, err := a.vdr.Read("did:indy:M9Z1siZMhVTdR3pH5zsxWm")
	//if err != nil {
	//	a.logger.Fatal(err)
	//}
	//
	//a.logger.Info("did doc fetched", doc.DIDDocument)
}
