package transport

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/agent/agent"
	"github.com/YasiruR/agent/transport/requests"
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

	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) handleCreateInvitation(w http.ResponseWriter, r *http.Request) {
	inv, err := s.agent.CreateInvitation()
	if err != nil {
		s.logger.Error(fmt.Sprintf(`create invitation - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(inv)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`encoding create invitation response - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleAcceptInvitation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var req requests.AcceptInv
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	res, err := s.agent.AcceptInvitation(req.Invitation)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`accept invitation - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
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
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`writing response - %v`, err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
