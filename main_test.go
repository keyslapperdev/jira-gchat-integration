package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("Registered Routes", func(t *testing.T) {
		routeNames := []string{
			"data",
			"healthCheck",
		}

		for _, name := range routeNames {
			route := getRouter().Get(name)
			assert.NotNilf(t, route, "No router found with name %q", name)
		}

	})

	t.Run("Data receive route", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(receiveData))
		payload := `{"text": "some text"}`

		resp, err := http.Post(server.URL, "content-type: application/json", bytes.NewBufferString(payload))
		assert.NoError(t, err, "Error posting JSON")

		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, payload, string(body), "Bad message returned")
	})

	t.Run("Health Check Route", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(healthCheck))

		resp, err := http.Get(server.URL)
		assert.NoError(t, err, "Error posting JSON")

		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, `{}`, string(body), "Bad message returned")
	})
}
