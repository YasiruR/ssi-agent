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
)

// agent endpoints
const (
	endpointCreateInv = `/out-of-band/create-invitation`
	endpointAcceptInv = `/out-of-band/receive-invitation`
	endpointGetConn   = `/connections/`
)

type Agent struct {
	port     int
	adminUrl string
	client   *http.Client
	logger   log.Logger
}

func New(port int, adminUrl string, logger log.Logger) *Agent {
	return &Agent{
		port:     port,
		adminUrl: adminUrl,
		client:   &http.Client{},
		logger:   logger,
	}
}

// CreateInvitation creates an invitation corresponding to out-of-band protocol
func (a *Agent) CreateInvitation() (*responses.CreateInvitation, error) {
	body := requests.CreateInvitation{
		Alias:              fmt.Sprintf("agent %d", a.port),
		HandshakeProtocols: []string{"did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/didexchange/1.0"},
		MyLabel:            "invitation to peer agent",
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
		var errRes responses.Error
		err = json.Unmarshal(data, &errRes)
		if err != nil {
			return nil, fmt.Errorf("response error encoding - %d", res.StatusCode)
		}
		return nil, fmt.Errorf("response error - %d [err = %v]", res.StatusCode, errRes)
	}

	var inv responses.CreateInvitation
	err = json.Unmarshal(data, &inv)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response - %v [%s]", err, string(data))
	}

	a.logger.Debug("invitation created for out-of-band protocol")
	return &inv, nil
}

// AcceptInvitation sends the received invitation to agent component and returns the response to sender (inviter)
func (a *Agent) AcceptInvitation(inv domain.Invitation) (response []byte, err error) {
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

	var accInv responses.AcceptInvitation
	err = json.Unmarshal(data, &accInv)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response - %v [%s]", err, string(data))
	}

	a.logger.Debug("invitation accepted for out-of-band protocol", accInv)
	return data, nil
}

func (a *Agent) Connection(connID string) (response []byte, err error) {
	res, err := a.client.Get(a.adminUrl + endpointGetConn + connID)
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
