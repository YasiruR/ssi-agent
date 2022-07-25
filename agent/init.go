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

type Agent struct {
	port     int
	adminUrl string
	client   *http.Client
	logger   log.Logger
}

func New() *Agent {
	return &Agent{}
}

// CreateInvitation creates an invitation corresponding to out-of-band protocol
func (a *Agent) CreateInvitation() (*responses.Invitation, error) {
	body := requests.CreateInvitation{
		Alias:              fmt.Sprintf("agent %d", a.port),
		HandshakeProtocols: []string{`did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/didexchange/1.0`},
		MyLabel:            "invitation to peer agent",
	}

	data, err := json.Marshal(&body)
	if err != nil {
		return nil, fmt.Errorf("request payload - %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, a.adminUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("http request - %v", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transport error - %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error - %d", res.StatusCode)
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body - %v", err)
	}

	var inv responses.Invitation
	err = json.Unmarshal(data, &inv)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response - %v", err)
	}

	return &inv, nil
}
