package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"io/ioutil"
	"log"
	"net/http"
)

type transport struct {
	router *mux.Router
	agent  *agent
}

func initHttpClient(agent *agent) {
	t := &transport{
		router: mux.NewRouter(),
		agent:  agent,
	}

	t.router.HandleFunc(`/connection`, t.handleCreateInvitation).Methods(http.MethodGet)
	t.router.HandleFunc(`/connection`, t.handleHandleInvitation).Methods(http.MethodPost)

	if err := http.ListenAndServe(":5005", t.router); err != nil {
		log.Fatal(err)
	}
}

func (t *transport) handleCreateInvitation(w http.ResponseWriter, r *http.Request) {
	inv, err := t.agent.CreateInvitation()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("agent sending invitation ", *inv)

	err = json.NewEncoder(w).Encode(inv)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *transport) handleHandleInvitation(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	inv := &didexchange.Invitation{}
	err = json.Unmarshal(data, inv)
	if err != nil {
		log.Fatal("unmarshall error: ", err)
	}

	id, err := t.agent.HandleInvitation(inv)
	if err != nil {
		log.Fatal("agent error: ", err)
	}

	fmt.Println("user connecting to ", id)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
}
