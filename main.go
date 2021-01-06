package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const port = ":7070"

var logger = StartLogger()

func main() {
	var (
		certFile    = os.Getenv("JIRABOT_CERT_FILE")
		certKeyFile = os.Getenv("JIRABOT_CERT_KEY_FILE")
		useSSL      = false
	)

	if os.Getenv("JIRABOT_USE_SSL") == "true" {
		useSSL = true
	}

	jira := JiraService{}
	jira.Authorize()

	chat := ChatService{}
	chat.Authorize()

	r := getRouter(jira, chat)

	fmt.Println("Running on port " + port)

	if useSSL {
		logger.Fatal(http.ListenAndServeTLS(port, certFile, certKeyFile, r))
	} else {
		logger.Fatal(http.ListenAndServe(port, r))
	}
}

// ChatPayload consumes the essential information from the data sent by
// google chat
type ChatPayload struct {
	Type    string  `json:"type"`
	Message Message `json:"message"`
}

// Message holds the arguments passed to the bot
type Message struct {
	Args string `json:"argumentText"`
}

// cleanup takes the provided args and pulls out the requied information
func cleanup(cp *ChatPayload) {
	cp.Message.Args = strings.Fields(cp.Message.Args)[0]

	re := regexp.MustCompile(`\D$`)
	for re.Match([]byte(cp.Message.Args)) {
		cp.Message.Args = string(cp.Message.Args[:len(cp.Message.Args)-1])
	}
}
