package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/YasiruR/agent/agent/requests"
	"github.com/YasiruR/agent/agent/responses"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
)

// agent endpoints
const (
	endpointCreateInv = `/out-of-band/create-invitation`
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
func (a *Agent) CreateInvitation() (*responses.Invitation, error) {
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

	var inv responses.Invitation
	err = json.Unmarshal(data, &inv)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response - %v [%s]", err, string(data))
	}

	a.logger.Debug("invitation created for out-of-band protocol")
	return &inv, nil
}
