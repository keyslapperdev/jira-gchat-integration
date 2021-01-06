package main

import (
	"os"

	jira "github.com/andygrunwald/go-jira"
)

var jiraAuthUser = os.Getenv("JIRABOT_AUTH_USER")
var jiraAuthPass = os.Getenv("JIRABOT_AUTH_PASS")

var jiraBaseURL = os.Getenv("JIRABOT_BASE_URL")

// GetAuthdJiraClient takes env data to authorize a session
// for jira -- in the future I plan to not use basic auth, but
// the oauth flow
func GetAuthdJiraClient() *jira.Client {
	logger.Trace("Authenticating Jira client")

	ac := jira.BasicAuthTransport{
		Username: jiraAuthUser,
		Password: jiraAuthPass,
	}

	client, err := jira.NewClient(ac.Client(), jiraBaseURL)
	if err != nil {
		logger.Fatal("Unable to connect to jira: " + err.Error())
	}

	return client
}
