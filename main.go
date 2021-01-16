package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

const port = ":7070"

var logger *logrus.Logger

func init() {
	writers := make([]io.Writer, 0)

	var logFile = os.Getenv("JIRABOT_LOG_FILE")
	if logFile != "" {
		fd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			log.Printf("Coudln't open log file: %v\nLogging to STDERR.", err)
		} else {
			writers = append(writers, fd)
		}
	} else {
		writers = append(writers, os.Stderr)
	}

	if os.Getenv("JIRABOT_LOG_TO_CHAT") == "true" {
		chatLW := ChatLogWriter{}
		chatLW.Authorize()

		writers = append(writers, chatLW)
	}

	logger = StartLogger(writers...)
}

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
