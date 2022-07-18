package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

type httpClient struct {
	*agent
	*mux.Router
	logger log.Logger
}

func newHttpClient(port int, agent *agent, logger log.Logger) {
	h := &httpClient{agent: agent, Router: mux.NewRouter(), logger: logger}

	h.Router.HandleFunc(`/invitation/create`, h.handleCreateInv).Methods(http.MethodPut)
	h.Router.HandleFunc(`/invitation/accept`, h.handleConnect).Methods(http.MethodPost)

	h.Router.HandleFunc(`/connection`, h.handleGetConn).Methods(http.MethodGet)
	h.Router.HandleFunc(`/send-offer`, h.handleSendOffer).Methods(http.MethodPost)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), h.Router); err != nil {
		h.logger.Fatal(err)
	}
}

func (h *httpClient) handleCreateInv(w http.ResponseWriter, _ *http.Request) {
	inv, err := h.agent.createInv()
	if err != nil {
		h.logger.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(inv)
	if err != nil {
		h.logger.Fatal(err)
	}
}

func (h *httpClient) handleConnect(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Fatal(err)
	}
	defer r.Body.Close()

	var req didexchange.Invitation
	err = json.Unmarshal(data, &req)
	if err != nil {
		h.logger.Fatal(err)
	}

	conn, err := h.agent.connect(&req)
	if err != nil {
		h.logger.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(conn)
	if err != nil {
		h.logger.Fatal(err)
	}
}

func (h *httpClient) handleGetConn(w http.ResponseWriter, r *http.Request) {
	conns, err := h.agent.getConn()
	if err != nil {
		h.logger.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(conns)
	if err != nil {
		h.logger.Fatal(err)
	}
}

func (h *httpClient) handleSendOffer(w http.ResponseWriter, r *http.Request) {
	h.agent.sendOffer()
}
