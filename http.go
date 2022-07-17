package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
)

type transport struct {
	router *mux.Router
	agent  *agent
	store  *store
	logger log.Logger
}

func initHttpClient(port string, agent *agent, store *store, logger log.Logger) {
	t := &transport{
		router: mux.NewRouter(),
		agent:  agent,
		store:  store,
		logger: logger,
	}

	t.router.HandleFunc(`/invitation/create`, t.handleCreateInvitation).Methods(http.MethodPut)
	t.router.HandleFunc(`/invitation/accept`, t.handleAcceptInvitation).Methods(http.MethodPost)

	// service endpoint
	t.router.HandleFunc(`/invitation/service`, t.handleServiceEndpoint).Methods(http.MethodGet)

	if err := http.ListenAndServe(":"+port, t.router); err != nil {
		t.logger.Fatal(err)
	}
}

func (t *transport) handleCreateInvitation(w http.ResponseWriter, _ *http.Request) {
	inv, err := t.agent.createInv()
	if err != nil {
		t.logger.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(inv)
	if err != nil {
		t.logger.Fatal(err)
	}
}

func (t *transport) handleAcceptInvitation(_ http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.logger.Fatal(err)
	}

	var req outofband.Invitation
	err = json.Unmarshal(data, &req)
	if err != nil {
		t.logger.Fatal(err)
	}

	_, err = t.agent.acceptInv(&req)
	if err != nil {
		t.logger.Fatal(err)
	}
}

//func (t *transport) handleGetConn(w http.ResponseWriter, r *http.Request) {
//	err := t.store.getConn()
//	if err != nil {
//		t.logger.Fatal(err, t.store.connID)
//	}
//}

func (t *transport) handleServiceEndpoint(w http.ResponseWriter, r *http.Request) {
	t.logger.Info("service endpoint called")
}
