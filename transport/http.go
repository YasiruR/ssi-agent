package transport

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/agent/agent"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
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
