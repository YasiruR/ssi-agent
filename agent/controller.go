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
	endpointSchemas   = `/schemas`
	endpointCredDef   = `/credential-definitions`
	endpointSendOffer = `/issue-credential-2.0/send-offer`
	endpointRecords   = `/issue-credential-2.0/records/`
)

type Agent struct {
	name     string
	adminUrl string
	client   *http.Client
	logger   log.Logger
	conns    *sync.Map // peer label to own connection ID map
	credMap  *sync.Map // agent label to credential exchange ID map
}

func New(name string, adminUrl string, logger log.Logger) *Agent {
	return &Agent{
		name:     name,
		adminUrl: adminUrl,
		client:   &http.Client{},
		logger:   logger,
		conns:    &sync.Map{},
		credMap:  &sync.Map{},
	}
}

func (a *Agent) AddConnection(label, connID string) {
	a.conns.Store(label, connID)
}

func (a *Agent) AddCredentialRecord(label, credExID string) {
	//var credIDs []string
	//val, ok := a.credMap.Load(label)
	//if ok {
	//	credIDs, ok = val.([]string)
	//	if !ok {
	//		a.logger.Error(fmt.Sprintf(`incompatible credential exchange ID list found for label %s`, label), val)
	//		return
	//	}
	//}
	//credIDs = append(credIDs, credExID)
	if a.name != label {
		a.credMap.Store(label, credExID)
		a.logger.Debug("credential record saved", label, credExID)
	}
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error - %d", res.StatusCode)
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body - %v", err)
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error - %d", res.StatusCode)
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body - %v", err)
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("accept response error - %d", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("accept reading body - %v", err)
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

	return a.post(a.adminUrl+endpointConn+connID+`/accept-request`, nil, fmt.Sprintf("connnection request accepted for id %s", connID))
}

// Connection fetches connection details for the given ID from agent endpoint and returns the response
func (a *Agent) Connection(connID string) (response []byte, err error) {
	return a.get(a.adminUrl+endpointConn+connID, fmt.Sprintf("connection fetched %s", connID))
}

// CreateSchema forwards the received schema directly to the agent (needs to be a Trust Anchor)
func (a *Agent) CreateSchema(schema []byte) (response []byte, err error) {
	return a.post(a.adminUrl+endpointSchemas, schema, "schema created")
}

// CreateCredentialDef forwards the received credential definition body directly to the agent to persist on ledger
func (a *Agent) CreateCredentialDef(def []byte) (response []byte, err error) {
	return a.post(a.adminUrl+endpointCredDef, def, "credential definition created")
}

// SendOffer takes domain.CredentialPreview and domain.IndySchemaMeta along with the recipient label which will then be
// used to send a credential offer to the (to-be) holder
func (a *Agent) SendOffer(cp domain.CredentialPreview, indySchema domain.IndySchemaMeta, to string) (response []byte, err error) {
	req := requests.Offer{}
	val, ok := a.conns.Load(to)
	if !ok {
		return nil, fmt.Errorf(`no connection found for the recipient %s`, to)
	}

	connID, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf(`connection ID corresponding to the recipient is not a string`)
	}

	req.ConnectionID = connID
	req.CredentialPreview = cp
	req.Filter.Indy = indySchema
	req.Comment = a.name
	a.logger.Debug("credential offer constructed", req)

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(`marshal error - %v`, err)
	}

	return a.post(a.adminUrl+endpointSendOffer, data, fmt.Sprintf("offer sent to %s", to))
}

// CredentialRecord finds the corresponding credential exchange ID from the in-memory map and fetches the credential record from the ledger
func (a *Agent) CredentialRecord(label string) (response []byte, err error) {
	val, ok := a.credMap.Load(label)
	if !ok {
		return nil, fmt.Errorf(`no credential offers received by %s`, label)
	}

	credExID, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf(`incompatible credential exchange ID found for label %s [%v]`, label, val)
	}

	return a.get(a.adminUrl+endpointRecords+credExID, fmt.Sprintf("credential record fetched with id %s", credExID))
}

// RequestCredential proceeds with requesting the credential from the issuer which corresponds to the given credential exchange ID
func (a *Agent) RequestCredential(credExID string) (response []byte, err error) {
	return a.post(a.adminUrl+endpointRecords+credExID+`/send-request`, nil, fmt.Sprintf("requested credential with id %s", credExID))
}

// IssueCredential proceeds with issuing the credential to the holder via connected agent
func (a *Agent) IssueCredential(credExID string) (response []byte, err error) {
	body := `{"comment": "issuing credential"}`
	return a.post(a.adminUrl+endpointRecords+credExID+`/issue`, []byte(body), fmt.Sprintf("issued credential with id %s", credExID))
}

// StoreCredential fetches the credential record by the given ID and stores it in the wallet of the holder. If user needs to fetch
// this stored credential directly from the wallet, id corresponding to `cred_id_stored` parameter of this response should be used.
func (a *Agent) StoreCredential(credExID string) (response []byte, err error) {
	return a.post(a.adminUrl+endpointRecords+credExID+`/store`, nil, fmt.Sprintf("stored credential with id %s", credExID))
}

// post proceeds with sending POST request
func (a *Agent) post(url string, body []byte, successLog string) (response []byte, err error) {
	res, err := a.client.Post(url, `application/json`, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf(`transport error - %v`, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error - %d", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body - %v", err)
	}

	a.logger.Debug(successLog)
	return data, nil
}

// get proceeds with sending GET request
func (a *Agent) get(url string, successLog string) (response []byte, err error) {
	res, err := a.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf(`transport error - %v`, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error - %d", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body - %v", err)
	}

	a.logger.Debug(successLog)
	return data, nil
}
