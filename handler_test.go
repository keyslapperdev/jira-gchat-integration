package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	chat "google.golang.org/api/chat/v1"
)

func TestServer(t *testing.T) {
	t.Run("Registered Routes", func(t *testing.T) {
		mjs := MockJiraService{}
		mcs := MockChatService{}

		routeNames := []string{
			"data",
			"healthCheck",
		}

		for _, name := range routeNames {
			route := getRouter(mjs, mcs).Get(name)
			assert.NotNilf(t, route, "No router found with name %q", name)
		}
	})

	t.Run("Data receive route", func(t *testing.T) {
		called1 := 0
		called2 := 0
		mjs := MockJiraService{&called1}
		mcs := MockChatService{&called2}

		server := httptest.NewServer(http.HandlerFunc(getDataHandler(mjs, mcs)))
		payload := `{"type": "MESSAGE", "message": { "text": "@Jira xxx-1234" }}`

		_, err := http.Post(server.URL, "application/json", bytes.NewBufferString(payload))
		assert.NoError(t, err, "Error posting JSON")
		assert.Equal(t, 1, *mjs.called, "Jira not called")
		assert.Equal(t, 1, *mcs.called, "Chat not called")
	})

	t.Run("Health Check Route", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(healthCheck))

		resp, err := http.Get(server.URL)
		assert.NoError(t, err, "Error getting healthcheck route")

		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, `{}`, string(body), "Bad message returned")
	})
}

type MockJiraService struct{ called *int }

func (mjs MockJiraService) GetTicketData(string) (*jira.Issue, error) {
	*mjs.called++
	return nil, nil
}

type MockChatService struct{ called *int }

func (mcs MockChatService) CreateIssueCard(*jira.Issue) (*chat.Message, error) {
	*mcs.called++
	return nil, nil
}
