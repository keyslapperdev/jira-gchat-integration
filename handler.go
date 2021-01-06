package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getRouter(jira JiraWorker, chat ChatWorker) *mux.Router {
	logger.Trace("instantiating router")

	r := mux.NewRouter()

	r.NewRoute().
		Name("data").
		Path("/data").
		Methods(http.MethodPost).
		HandlerFunc(getDataHandler(jira, chat))
	r.NewRoute().
		Name("healthCheck").
		Path("/health").
		Methods(http.MethodGet).
		HandlerFunc(healthCheck)

	return r
}

func getDataHandler(jira JiraWorker, chat ChatWorker) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		logger.Trace("hit on data route")

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

		message, err := chat.CreateIssueCard(tData)
		if err != nil {
			http.Error(rw, "Error creating card: "+err.Error(), http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(message)
		return
	}
}

func healthCheck(rw http.ResponseWriter, r *http.Request) {
	logger.Trace("hit on health route")

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("{}"))
}
