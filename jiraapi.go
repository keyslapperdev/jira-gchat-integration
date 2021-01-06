package main

import (
	"github.com/andygrunwald/go-jira"
)

// JiraWorker words
type JiraWorker interface {
	GetTicketData(ChatPayload) (*jira.Issue, error)
}

// JiraService is a wrapper for jira.Client
type JiraService struct {
	*jira.Client
}

// Authorize is a method on JiraService to authorize the client
// It's made in this way so that the authorization can happen only
// once
func (j *JiraService) Authorize() {
	j.Client = getAuthdJiraClient()
}

// GetTicketData fetches data from JIRA given the id.
func (j JiraService) GetTicketData(payload ChatPayload) (*jira.Issue, error) {
	issue, _, err := j.Issue.Get(payload.Message.Args, nil)
	if err != nil {
		return nil, err
	}

	return issue, nil
}
