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

	s.logger.Info(fmt.Sprintf("webhook server started listening on %d", s.port))
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var req requests.Connections
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	s.logger.Debug("webhook received for connection", req)
}
