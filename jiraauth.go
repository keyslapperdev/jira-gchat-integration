package main

import (
	"os"

	jira "github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
)

var jiraAuthUser = os.Getenv("JIRABOT_AUTH_USER")
var jiraAuthPass = os.Getenv("JIRABOT_AUTH_PASS")

var jiraBaseURL = os.Getenv("JIRABOT_BASE_URL")

func getAuthdJiraClient() *jira.Client {
	ac := jira.BasicAuthTransport{
		Username: jiraAuthUser,
		Password: jiraAuthPass,
	}

	client, err := jira.NewClient(ac.Client(), jiraBaseURL)
	if err != nil {
		log.Fatal("Unable to connect to jira: " + err.Error())
	}

	return client
}
