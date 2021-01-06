package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("Registered Routes", func(t *testing.T) {
		mjs := MockJiraService{}

		routeNames := []string{
			"data",
			"healthCheck",
		}

		for _, name := range routeNames {
			route := getRouter(mjs).Get(name)
			assert.NotNilf(t, route, "No router found with name %q", name)
		}
	})

	t.Run("Data receive route", func(t *testing.T) {
		count := 0
		mjs := MockJiraService{&count}

		server := httptest.NewServer(http.HandlerFunc(getDataHandler(mjs)))
		payload := `{"type": "MESSAGE", "message": { "argumentText": "xxx-1234" }}`

		_, err := http.Post(server.URL, "application/json", bytes.NewBufferString(payload))
		assert.NoError(t, err, "Error posting JSON")
		assert.Equal(t, 1, *mjs.called, "Jira not called")
	})

	t.Run("Health Check Route", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(healthCheck))

		resp, err := http.Get(server.URL)
		assert.NoError(t, err, "Error getting healthcheck route")

		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, `{}`, string(body), "Bad message returned")
	})
}

type MockJiraService struct {
	called *int
}

func (mjs MockJiraService) GetTicketData(p ChatPayload) (*jira.Issue, error) {
	*mjs.called++

	return nil, nil
}
