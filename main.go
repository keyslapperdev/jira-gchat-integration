package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const port = ":7070"

func init() { log.SetReportCaller(true) }

func main() {
	var config = loadConfig()

	var (
		certFile    = config.certFile
		certKeyFile = config.certKeyFile
		useSSL      = config.useSSL
	)

	r := getRouter()

	fmt.Println("Running on port " + port)

	if useSSL {
		log.Fatal(http.ListenAndServeTLS(port, certFile, certKeyFile, r))
	} else {
		log.Fatal(http.ListenAndServe(port, r))
	}
}
