package agent

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/agent/agent"
	"github.com/YasiruR/agent/transport/agent/requests"
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
	s.router.HandleFunc(`/invitation/create`, s.handleCreateInvitation).Methods(http.MethodPost)
	s.router.HandleFunc(`/invitation/accept`, s.handleAcceptInvitation).Methods(http.MethodPost)

	s.router.HandleFunc(`/connection/{id}`, s.handleGetConnection).Methods(http.MethodGet)
	s.router.HandleFunc(`/connection/accept-request/{their_label}`, s.handleAcceptRequest).Methods(http.MethodPost)

	s.router.HandleFunc(`/schema/create`, s.handleCreateSchema).Methods(http.MethodPost)

	s.logger.Info(fmt.Sprintf("controller started listening on %d", s.port))
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) handleCreateInvitation(w http.ResponseWriter, _ *http.Request) {
	res, err := s.agent.CreateInvitation()
	if err != nil {
		s.logger.Error(fmt.Sprintf(`create invitation - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`writing response - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleAcceptInvitation(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var req requests.AcceptInvitation
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := s.agent.AcceptInvitation(req.Invitation)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`accept invitation - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`writing response - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleAcceptRequest(w http.ResponseWriter, r *http.Request) {
	label := mux.Vars(r)[`their_label`]
	res, err := s.agent.AcceptRequest(label)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`accept request - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`writing response - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleGetConnection(w http.ResponseWriter, r *http.Request) {
	connID := mux.Vars(r)[`id`]
	res, err := s.agent.Connection(connID)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`get connection - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`writing response - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleCreateSchema(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	res, err := s.agent.CreateSchema(data)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`create schema - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`writing response - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
