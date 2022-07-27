package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/YasiruR/agent/agent/requests"
	"github.com/YasiruR/agent/agent/responses"
	"github.com/YasiruR/agent/domain"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"sync"
)

// agent endpoints
const (
	endpointCreateInv = `/connections/create-invitation`
	endpointAcceptInv = `/connections/receive-invitation`
	endpointConn      = `/connections/`
)

type Agent struct {
	name     string
	adminUrl string
	client   *http.Client
	logger   log.Logger
	conns    *sync.Map // peer label to own connection ID map
}

func New(name string, adminUrl string, logger log.Logger) *Agent {
	return &Agent{
		name:     name,
		adminUrl: adminUrl,
		client:   &http.Client{},
		logger:   logger,
		conns:    &sync.Map{},
	}
}

func (a *Agent) AddConnection(label, connID string) {
	a.conns.Store(label, connID)
}

// CreateInvitation creates an invitation corresponding to out-of-band protocol
func (a *Agent) CreateInvitation() (response []byte, err error) {
	body := requests.CreateInvitation{
		Alias:              fmt.Sprintf("agent %s", a.name),
		HandshakeProtocols: []string{"did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/didexchange/1.0"},
		MyLabel:            a.name,
		UsePublicDid:       false,
	}

	data, err := json.Marshal(&body)
	if err != nil {
		return nil, fmt.Errorf("request payload - %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, a.adminUrl+endpointCreateInv, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("http request - %v", err)
	}

	res, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transport error - %v", err)
	}
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body - %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error - %d", res.StatusCode)
	}

	var inv responses.CreateInvitation
	err = json.Unmarshal(data, &inv)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response - %v [%s]", err, string(data))
	}

	a.logger.Debug("invitation created for did-exchange protocol", inv)
	return data, nil
}

// AcceptInvitation sends the received invitation to agent component for storage. If successful, controller proceeds with
// accepting the invitation with the connection id and returns the response to sender (inviter)
func (a *Agent) AcceptInvitation(inv domain.Invitation) (response []byte, err error) {
	recInv, err := a.receiveInvitation(inv)
	if err != nil {
		return nil, fmt.Errorf(`receive invitation - %v`, err)
	}

	return a.acceptInvitation(recInv.ConnectionID)
}

func (a *Agent) receiveInvitation(inv domain.Invitation) (*responses.ReceiveInvitation, error) {
	data, err := json.Marshal(&inv)
	if err != nil {
		return nil, fmt.Errorf("request payload - %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, a.adminUrl+endpointAcceptInv, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("http request - %v", err)
	}
	req.Header.Add(`accept`, `application/json`)
	req.Header.Add(`Content-Type`, `application/json`)

	res, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transport error - %v", err)
	}
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body - %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error - %d", res.StatusCode)
	}

	var recInv responses.ReceiveInvitation
	err = json.Unmarshal(data, &recInv)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response - %v [%s]", err, string(data))
	}

	a.logger.Debug("invitation received for did-exchange protocol", recInv)
	return &recInv, nil
}

func (a *Agent) acceptInvitation(connID string) (response []byte, err error) {
	req, err := http.NewRequest(http.MethodPost, a.adminUrl+endpointConn+connID+`/accept-invitation`, nil)
	if err != nil {
		return nil, fmt.Errorf(`request error - %v`, err)
	}

	// add label to the connection
	params := req.URL.Query()
	params.Add(`my_label`, a.name)
	req.URL.RawQuery = params.Encode()

	res, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("accept transport error - %v", err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("accept reading body - %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("accept response error - %d", res.StatusCode)
	}

	a.logger.Debug("invitation accepted for did-exchange protocol")
	return data, nil
}

// AcceptRequest maps the label to connection ID received by webhook and proceeds with accepting connection request via agent
func (a *Agent) AcceptRequest(label string) (response []byte, err error) {
	val, ok := a.conns.Load(label)
	if !ok {
		return nil, fmt.Errorf(`no connection found for the label`)
	}

	connID, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf(`connection ID corresponding to the label is not a string`)
	}

	res, err := a.client.Post(a.adminUrl+endpointConn+connID+`/accept-request`, `application/json`, nil)
	if err != nil {
		return nil, fmt.Errorf("accept transport error - %v", err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("accept reading body - %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("accept response error - %d", res.StatusCode)
	}

	a.logger.Debug("connection request accepted", connID)
	return data, nil
}

func (a *Agent) Connection(connID string) (response []byte, err error) {
	res, err := a.client.Get(a.adminUrl + endpointConn + connID)
	if err != nil {
		return nil, fmt.Errorf("transport error - %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error - %d", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response - %v", err)
	}

	var conn domain.Connection
	err = json.Unmarshal(data, &conn)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error - %v", err)
	}

	a.logger.Debug("connection fetched", conn)
	return data, nil
}
