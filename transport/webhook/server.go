package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/agent/agent"
	"github.com/YasiruR/agent/transport/webhook/requests"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Server struct {
	port   int
	router *mux.Router
	agent  *agent.Agent
	logger log.Logger
}

func New(port int, agent *agent.Agent, logger log.Logger) *Server {
	return &Server{port: port, router: mux.NewRouter(), agent: agent, logger: logger}
}

func (s *Server) Serve() {
	s.router.HandleFunc(`/topic/connections/`, s.handleConnections).Methods(http.MethodPost)
	s.router.HandleFunc(`/topic/issue_credential_v2_0/`, s.handleCredentials).Methods(http.MethodPost)
	s.router.HandleFunc(`/topic/issue_credential_v2_0_indy/`, s.handleIndyCredentials).Methods(http.MethodPost)
	s.router.HandleFunc(`/topic/present_proof_v2_0/`, s.handlePresentProof).Methods(http.MethodPost)

	s.logger.Info(fmt.Sprintf("webhook server started listening on %d", s.port))
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) handleConnections(_ http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		return
	}

	var req requests.Connections
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.logger.Error(err)
		return
	}

	s.logger.Debug("webhook received for connection", req)
	s.agent.AddConnection(req.TheirLabel, req.ConnectionID)
}

func (s *Server) handleCredentials(_ http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		return
	}

	var req requests.IssueCredentials
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.logger.Error(err)
		return
	}

	s.logger.Debug("webhook received for credentials", req)

	// workaround to proceed with credential offers
	if req.CredIssue.ID == `` && req.CredOffer.ID != `` {
		s.agent.AddCredentialRecord(req.CredOffer.Comment, req.CredExID)
	}
}

func (s *Server) handleIndyCredentials(_ http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		return
	}

	var req requests.IssueCredentialsIndy
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.logger.Error(err)
		return
	}

	s.logger.Debug("webhook received for credentials for indy", req)
}

func (s *Server) handlePresentProof(_ http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		return
	}

	var req requests.PresentationProof
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.logger.Error(err)
		return
	}

	s.logger.Debug("webhook received for proof presentation", req)
	s.agent.AddPresentationRecord(req.PresRequest.Comment, req.PresExID, req.ByFormat.PresRequest)
}
