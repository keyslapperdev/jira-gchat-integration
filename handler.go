package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// Maintainer would be the name of ... me, it's just me. :(
const Maintainer = "Alexander Wilcots (alexander.wilcots@endurance.com)"

// BaseRoute contains the base route for the bot, i.e. /api/v1/
var BaseRoute = os.Getenv("BASE_ROUTE")

func getRouter(jira JiraWorker, chat ChatWorker) *mux.Router {
	logger.Trace("instantiating router")

	r := mux.NewRouter()
	s := r.PathPrefix(BaseRoute).Subrouter()

	s.NewRoute().
		Name("data").
		Path("/data").
		Methods(http.MethodPost).
		HandlerFunc(getDataHandler(jira, chat))
	s.NewRoute().
		Name("healthCheck").
		Path("/health").
		Methods(http.MethodGet).
		HandlerFunc(healthCheck)

	return r
}

func getDataHandler(jira JiraWorker, chat ChatWorker) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		logger.Trace("hit on data route")

		payload := ChatPayload{}
		json.NewDecoder(r.Body).Decode(&payload)

		tID := getTicketID(payload)

		if tID == "" {
			logger.Infof("bad data provided: %s", payload.Message.Text)
			http.Error(rw, fmt.Sprintf(`{"text": "[!] Please enter a vaild Jira ticket id %s"}`, usage()), http.StatusOK)

			return
		}

		tData, err := jira.GetTicketData(tID)
		if err != nil {
			logger.Error("Error with Jira: " + err.Error())

			if strings.Contains(err.Error(), "permission") {
				http.Error(rw, fmt.Sprintf(`{"text": "My apologies, my jira user (svcjirahgeng) doesn't have access to view this ticket (%s).\nIf possible, please authorize me to view it better use out of me."}`, tID), http.StatusOK)
			} else if strings.Contains(err.Error(), "Not Exist") {
				http.Error(rw, fmt.Sprintf(`{"text": "The requested ticket %s does not seem to exist."}`, tID), http.StatusOK)
			} else {
				http.Error(rw, fmt.Sprintf(`{"text": "Jira: %s\nPlease contact %s with a paste of this error for assistance."}`, err.Error(), Maintainer), http.StatusOK)
			}

			return
		}

		message, err := chat.CreateIssueCard(tData)
		if err != nil {
			logger.Error("Error creating card: " + err.Error())
			http.Error(rw, fmt.Sprintf(`{"text": "Chat Card: %s\nPlease contact %s with a paste of this error for assistance."}`, err.Error(), Maintainer), http.StatusInternalServerError)
			return
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
