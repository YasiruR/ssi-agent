package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/YasiruR/agent/agent/models"
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
	endpointCreateInv    = `/connections/create-invitation`
	endpointAcceptInv    = `/connections/receive-invitation`
	endpointConn         = `/connections/`
	endpointSchemas      = `/schemas`
	endpointCredDef      = `/credential-definitions`
	endpointSendOffer    = `/issue-credential-2.0/send-offer`
	endpointSendCredAuto = `/issue-credential-2.0/send`
	endpointCredRecords  = `/issue-credential-2.0/records/`
	endpointSendProofReq = `/present-proof-2.0/send-request`
	endpointProofRecords = `/present-proof-2.0/records/`
	endpointCredentials  = `/credentials`
)

type Agent struct {
	name     string
	adminUrl string
	client   *http.Client
	logger   log.Logger
	connMap  *sync.Map // peer agent label to own connection ID map
	credMap  *sync.Map // peer agent label to credential exchange ID map
	proofMap *sync.Map // peer agent label to proof exchange ID map
}

func New(name string, adminUrl string, logger log.Logger) *Agent {
	return &Agent{
		name:     name,
		adminUrl: adminUrl,
		client:   &http.Client{},
		logger:   logger,
		connMap:  &sync.Map{},
		credMap:  &sync.Map{},
		proofMap: &sync.Map{},
	}
}

func (a *Agent) AddConnection(label, connID string) {
	a.connMap.Store(label, connID)
}

func (a *Agent) GetConnectionByLabel(label string) (string, error) {
	val, ok := a.connMap.Load(label)
	if !ok {
		return ``, fmt.Errorf(`no connection found for the recipient %s`, label)
	}

	connID, ok := val.(string)
	if !ok {
		return ``, fmt.Errorf(`connection ID corresponding to the recipient is not a string`)
	}

	return connID, nil
}

func (a *Agent) AddCredentialRecord(label, credExID string) {
	if a.name != label {
		a.credMap.Store(label, credExID)
		a.logger.Debug("credential record saved", label, credExID)
	}
}

func (a *Agent) AddPresentationRecord(label, presExID string, pr domain.PresentationRequest) {
	if a.name != label {
		a.proofMap.Store(label, models.ProofPresentation{PresExID: presExID, PresReq: pr})
		a.logger.Debug("proof record saved", label, presExID)
	}
}

func (a *Agent) GetPresentationRecord(label string) (models.ProofPresentation, error) {
	val, ok := a.proofMap.Load(label)
	if !ok {
		return models.ProofPresentation{}, fmt.Errorf(`presentation record does not exist for label %s`, label)
	}

	pp, ok := val.(models.ProofPresentation)
	if !ok {
		return models.ProofPresentation{}, fmt.Errorf(`incompatible proof record found for label %s [%v]`, label, val)
	}

	return pp, nil
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
	connID, err := a.GetConnectionByLabel(label)
	if err != nil {
		return nil, fmt.Errorf(`get connection by label - %v`, err)
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

// SendCredentialOffer takes domain.CredentialPreview and domain.IndySchemaMeta along with the recipient label which will then be
// used to send a credential offer to the (to-be) holder. Setting auto_remove of offer to true removes credential exchange
// record automatically after the protocol completes.
func (a *Agent) SendCredentialOffer(cp domain.CredentialPreview, indySchema domain.IndySchemaMeta, to string) (response []byte, err error) {
	req := requests.Offer{}
	req.CredentialPreview = cp
	req.Filter.Indy = indySchema
	req.Comment = a.name
	req.ConnectionID, err = a.GetConnectionByLabel(to)
	if err != nil {
		return nil, fmt.Errorf(`get connection by label - %v`, err)
	}
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

	return a.get(a.adminUrl+endpointCredRecords+credExID, fmt.Sprintf("credential record fetched with id %s", credExID))
}

// RequestCredential proceeds with requesting the credential from the issuer which corresponds to the given credential exchange ID
func (a *Agent) RequestCredential(credExID string) (response []byte, err error) {
	return a.post(a.adminUrl+endpointCredRecords+credExID+`/send-request`, nil, fmt.Sprintf("requested credential with id %s", credExID))
}

// IssueCredential proceeds with issuing the credential to the holder via connected agent
func (a *Agent) IssueCredential(credExID string) (response []byte, err error) {
	body := `{"comment": "issuing credential"}`
	return a.post(a.adminUrl+endpointCredRecords+credExID+`/issue`, []byte(body), fmt.Sprintf("issued credential with id %s", credExID))
}

// StoreCredential fetches the credential record by the given ID and stores it in the wallet of the holder. If user needs to fetch
// this stored credential directly from the wallet, id corresponding to `cred_id_stored` parameter of this response should be used.
func (a *Agent) StoreCredential(credExID string) (response []byte, err error) {
	return a.post(a.adminUrl+endpointCredRecords+credExID+`/store`, nil, fmt.Sprintf("stored credential with id %s", credExID))
}

// SendCredentialAuto starts from sending an offer for a credential and follows an automated process for the rest of the steps.
// This needs the holder to enable auto-responsiveness to credential offers.
func (a *Agent) SendCredentialAuto(cp domain.CredentialPreview, indySchema domain.IndySchemaMeta, to string) (response []byte, err error) {
	req := requests.Offer{}
	req.CredentialPreview = cp
	req.Filter.Indy = indySchema
	req.Comment = a.name
	req.ConnectionID, err = a.GetConnectionByLabel(to)
	if err != nil {
		return nil, fmt.Errorf(`get connection by label - %v`, err)
	}
	a.logger.Debug("credential offer constructed for automated process", req)

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(`marshal error - %v`, err)
	}

	return a.post(a.adminUrl+endpointSendCredAuto, data, fmt.Sprintf("automated process started by sending offer to %s", to))
}

// SendProofRequest finds the corresponding connection ID of the peer agent label and sends a proof request with the given body
func (a *Agent) SendProofRequest(pr domain.PresentationRequest, to string) (response []byte, err error) {
	connID, err := a.GetConnectionByLabel(to)
	if err != nil {
		return nil, fmt.Errorf(`get connection by label - %v`, err)
	}

	req := requests.ProofRequest{Comment: a.name, ConnectionID: connID, PresentReq: pr}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(`marshal error - %v`, err)
	}

	return a.post(a.adminUrl+endpointSendProofReq, data, fmt.Sprintf(`proof request received by %s`, to))
}

// PresentProof sends the presentation format of the proof to a verifier given by the peer agent label. It first fetches
// existing credentials from the holder's wallet and gets the corresponding record stored by the webhook. This is then
// used to extract requested attributes by the verifier along with the existing credentials from the holder's proof in order
// to construct the presentation proof.
func (a *Agent) PresentProof(to string) (response []byte, err error) {
	// todo only requested attributes atm
	creds, err := a.getCredentialsFromWallet()
	if err != nil {
		return nil, fmt.Errorf(`get credentials failed - %v`, err)
	}

	pp, err := a.GetPresentationRecord(to)
	if err != nil {
		return nil, fmt.Errorf(`get record failed - %v`, err)
	}
	proof := a.constructProof(pp.PresReq.Indy.RequestedAttributes, creds)

	return a.sendProofPresentation(pp.PresExID, proof)
}

// constructProof checks the values from the existing credentials in the wallet to match with requested attributes
func (a *Agent) constructProof(reqAttributes map[string]domain.Attribute, creds []responses.WalletCredential) requests.ProofPresentation {
	var proof requests.ProofPresentation
	proof.Indy.RequestedAttributes = make(map[string]requests.AdditionalProp)
	proof.Indy.RequestedPredicates = make(map[string]requests.AdditionalProp)
	proof.Indy.SelfAttestedAttributes = make(map[string]string)
attrLoop:
	for reqAttrKey, reqAttr := range reqAttributes {
		for _, cred := range creds {
			for credAttrName, _ := range cred.Attrs {
				// checks if stored attribute name and cred def ID are same as the requested
				if reqAttr.Name == credAttrName {
					for _, restrict := range reqAttr.Restrictions {
						if restrict.CredDefID == cred.CredDefID {
							proof.Indy.RequestedAttributes[reqAttrKey] = requests.AdditionalProp{CredID: cred.Referent, Revealed: true}
							continue attrLoop
						}
					}
				}
			}
		}
	}

	a.logger.Debug(`presentation proof constructed`, proof)
	return proof
}

func (a *Agent) sendProofPresentation(presExID string, proofPres requests.ProofPresentation) (response []byte, err error) {
	data, err := json.Marshal(proofPres)
	if err != nil {
		return nil, fmt.Errorf(`marhsal error - %v`, err)
	}

	return a.post(a.adminUrl+endpointProofRecords+presExID+`/send-presentation`, data, fmt.Sprintf(`presentation sent with exchange id %s`, presExID))
}

func (a *Agent) getCredentialsFromWallet() ([]responses.WalletCredential, error) {
	data, err := a.get(a.adminUrl+endpointCredentials, fmt.Sprintf(`fetched credentials from the wallet of %s`, a.name))
	if err != nil {
		return nil, err
	}

	var res responses.Credentials
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, fmt.Errorf(`unmarshal error - %v`, err)
	}

	return res.Results, nil
}

func (a *Agent) VerifyProof(presExID string) (response []byte, err error) {
	return a.post(a.adminUrl+endpointProofRecords+presExID+`/verify-presentation`, nil, fmt.Sprintf(`verified presentation proof %s`, presExID))
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
