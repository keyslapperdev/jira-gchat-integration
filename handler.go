package main

import (
	"encoding/json"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

func getRouter(jira JiraWorker) *mux.Router {
	r := mux.NewRouter()

	r.NewRoute().
		Name("data").
		Path("/data").
		Methods(http.MethodPost).
		HandlerFunc(getDataHandler(jira))
	r.NewRoute().
		Name("healthCheck").
		Path("/health").
		Methods(http.MethodGet).
		HandlerFunc(healthCheck)

	return r
}

func getDataHandler(jira JiraWorker) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		types, exist := r.Header["Content-Type"]
		if !exist || types[0] != "application/json" {
			http.Error(rw, "must set content type to application/json", http.StatusBadRequest)
			return
		}

		payload := ChatPayload{}
		json.NewDecoder(r.Body).Decode(&payload)

		cleanup(&payload)

		tData, err := jira.GetTicketData(payload)
		if err != nil {
			http.Error(rw, "Error with Jira: "+err.Error(), http.StatusInternalServerError)
		}

		spew.Dump(tData)

		//card := buildChatCard(jInfo)
		//json.NewEncoder(rw).Encode(card)

		return
	}
}

func healthCheck(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("{}"))
}