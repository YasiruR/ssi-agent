package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	outofbandv22 "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/outofbandv2"
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

	t.router.HandleFunc(`/connection/create-invitation`, t.handleCreateInvitation).Methods(http.MethodPut)
	t.router.HandleFunc(`/connection/accept-invitation`, t.handleAcceptInvitation).Methods(http.MethodPost)

	//t.router.HandleFunc(`/connection/{id}`, t.handleGetConnection).Methods(http.MethodGet)
	//t.router.HandleFunc(`/connection/handle-invitation`, t.handleHandleInvitation).Methods(http.MethodPost)
	//t.router.HandleFunc(`/send-offer/{id}`, t.handleSendOffer).Methods(http.MethodPut)

	if err := http.ListenAndServe(":"+port, t.router); err != nil {
		t.logger.Fatal(err)
	}
}

func (t *transport) handleCreateInvitation(w http.ResponseWriter, _ *http.Request) {
	inv, err := t.agent.createInvitation()
	if err != nil {
		t.logger.Fatal(err)
	}

	t.logger.Debug("agent sending invitation: ", inv)
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
	defer r.Body.Close()

	inv := outofbandv22.Invitation{}
	err = json.Unmarshal(data, &inv)
	if err != nil {
		t.logger.Fatal("unmarshall error: ", err)
	}

	t.logger.Debug("agent receiving invitation: ", inv)
	err = t.agent.acceptInvitation(&inv)
	if err != nil {
		t.logger.Fatal(err)
	}
}

//func (t *transport) handleGetConnection(w http.ResponseWriter, r *http.Request) {
//	params := mux.Vars(r)
//	conn, err := t.agent.connection(params[`id`])
//	if err != nil {
//		t.logger.Fatal(err)
//	}
//
//	w.WriteHeader(http.StatusOK)
//	_, err = w.Write([]byte(conn.State))
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		t.logger.Fatal(err)
//	}
//}

//func (t *transport) handleCreateInvitation(w http.ResponseWriter, _ *http.Request) {
//	inv, err := t.agent.createInvitation()
//	if err != nil {
//		t.logger.Fatal(err)
//	}
//
//	t.logger.Debug("agent sending invitation ")
//	fmt.Println("service: ", inv.ServiceEndpoint)
//	fmt.Println("recipient keys: ", inv.RecipientKeys)
//	fmt.Println("id: ", inv.ID)
//	fmt.Println("label: ", inv.Label)
//	fmt.Println("did: ", inv.DID)
//	fmt.Println("routing keys: ", inv.DID)
//	fmt.Println("type: ", inv.Type)
//
//	err = json.NewEncoder(w).Encode(inv)
//	if err != nil {
//		t.logger.Fatal(err)
//	}
//}

//func (t *transport) handleAcceptInvitation(w http.ResponseWriter, r *http.Request) {
//	data, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		t.logger.Fatal(err)
//	}
//	defer r.Body.Close()
//
//	inv := didexchange.Invitation{}
//	err = json.Unmarshal(data, &inv)
//	if err != nil {
//		t.logger.Fatal("unmarshall error: ", err)
//	}
//
//	t.logger.Debug("agent receiving invitation ")
//	fmt.Println("service: ", inv.ServiceEndpoint)
//	fmt.Println("recipient keys: ", inv.RecipientKeys)
//	fmt.Println("id: ", inv.ID)
//	fmt.Println("label: ", inv.Label)
//	fmt.Println("did: ", inv.DID)
//	fmt.Println("routing keys: ", inv.DID)
//	fmt.Println("type: ", inv.Type)
//
//	err = t.agent.acceptInvitation(&inv)
//	if err != nil {
//		t.logger.Fatal(err)
//	}
//}

//func (t *transport) handleHandleInvitation(w http.ResponseWriter, r *http.Request) {
//	data, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		t.logger.Fatal(err)
//	}
//	defer r.Body.Close()
//
//	inv := &didexchange.Invitation{}
//	err = json.Unmarshal(data, inv)
//	if err != nil {
//		t.logger.Fatal("unmarshall error: ", err)
//	}
//
//	id, err := t.agent.handleInvitation(inv)
//	if err != nil {
//		t.logger.Fatal("agent error: ", err)
//	}
//
//	fmt.Println("user connecting to ", id)
//
//	w.WriteHeader(http.StatusOK)
//	_, err = w.Write([]byte(id))
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		t.logger.Fatal(err)
//	}
//}

//func (t *transport) handleSendOffer(_ http.ResponseWriter, r *http.Request) {
//	params := mux.Vars(r)
//	conn, err := t.agent.connection(params[`id`])
//	if err != nil {
//		t.logger.Fatal(err)
//	}
//
//	fmt.Println("CONNECTION: ", conn.MyDID, conn.TheirDID)
//	fmt.Println("RECORD: ", conn.Record.MyDID, conn.Record.TheirDID)
//
//	err = t.agent.sendVCOffer(conn.Record)
//	if err != nil {
//		t.logger.Fatal(err)
//	}
//}
