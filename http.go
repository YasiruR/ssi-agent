package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
)

type transport struct {
	router *mux.Router
	agent  *agent
	logger log.Logger
}

func initHttpClient(port string, agent *agent, logger log.Logger) {
	t := &transport{
		router: mux.NewRouter(),
		agent:  agent,
		logger: logger,
	}

	t.router.HandleFunc(`/connection/{id}`, t.handleGetConnection).Methods(http.MethodGet)
	t.router.HandleFunc(`/connection/create-invitation`, t.handleCreateInvitation).Methods(http.MethodPut)
	t.router.HandleFunc(`/connection/handle-invitation`, t.handleHandleInvitation).Methods(http.MethodPost)

	if err := http.ListenAndServe(":"+port, t.router); err != nil {
		t.logger.Fatal(err)
	}
}

func (t *transport) handleGetConnection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	conn, err := t.agent.connection(params[`id`])
	if err != nil {
		t.logger.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(conn.State))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		t.logger.Fatal(err)
	}
}

func (t *transport) handleCreateInvitation(w http.ResponseWriter, r *http.Request) {
	inv, err := t.agent.createInvitation()
	if err != nil {
		t.logger.Fatal(err)
	}

	fmt.Println("agent sending invitation ", inv.Invitation)

	err = json.NewEncoder(w).Encode(inv)
	if err != nil {
		t.logger.Fatal(err)
	}
}

func (t *transport) handleHandleInvitation(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.logger.Fatal(err)
	}
	defer r.Body.Close()

	inv := &didexchange.Invitation{}
	err = json.Unmarshal(data, inv)
	if err != nil {
		t.logger.Fatal("unmarshall error: ", err)
	}

	id, err := t.agent.handleInvitation(inv)
	if err != nil {
		t.logger.Fatal("agent error: ", err)
	}

	fmt.Println("user connecting to ", id)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		t.logger.Fatal(err)
	}
}
