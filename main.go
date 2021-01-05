package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const port = ":7070"

var config = loadConfig()

var (
	certFile    = config.certFile
	certKeyFile = config.certKeyFile
)

func init() { log.SetReportCaller(true) }

func main() {
	r := getRouter()

	fmt.Println("Running on port " + port)

	if config.useSSL {
		log.Fatal(http.ListenAndServeTLS(port, certFile, certKeyFile, r))
	} else {
		log.Fatal(http.ListenAndServe(port, r))
	}
}

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

type botConfig struct {
	certFile, certKeyFile string
	useSSL                bool
}

func loadConfig() (config botConfig) {
	config = botConfig{
		certFile:    os.Getenv("JIRABOT_CERT_FILE"),
		certKeyFile: os.Getenv("JIRABOT_CERT_KEY_FILE"),
		useSSL:      false,
	}

	if os.Getenv("JIRABOT_USE_SSL") == "true" {
		config.useSSL = true
	}

	return
}
