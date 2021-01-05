package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func getRouter() *mux.Router {
	r := mux.NewRouter()

	r.NewRoute().
		Name("data").
		Path("/data").
		Methods(http.MethodPost).
		HandlerFunc(receiveData)
	r.NewRoute().
		Name("healthCheck").
		Path("/health").
		Methods(http.MethodGet).
		HandlerFunc(healthCheck)

	return r
}

func receiveData(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "unable to parse body", http.StatusBadRequest)
		log.Fatal("unable to parse body")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(body)
	return
}

func healthCheck(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("{}"))
}
