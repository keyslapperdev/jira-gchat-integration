package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

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

	jira := JiraService{}
	jira.Authorize()

	r := getRouter(jira)

	fmt.Println("Running on port " + port)

	if useSSL {
		log.Fatal(http.ListenAndServeTLS(port, certFile, certKeyFile, r))
	} else {
		log.Fatal(http.ListenAndServe(port, r))
	}
}

type ChatPayload struct {
	Type    string  `json:"type"`
	Message Message `json:"message"`
}

type Message struct {
	Args string `json:"argumentText"`
}

func cleanup(cp *ChatPayload) {
	cp.Message.Args = strings.Fields(cp.Message.Args)[0]

	re := regexp.MustCompile(`\D$`)
	for re.Match([]byte(cp.Message.Args)) {
		cp.Message.Args = string(cp.Message.Args[:len(cp.Message.Args)-1])
	}
}
