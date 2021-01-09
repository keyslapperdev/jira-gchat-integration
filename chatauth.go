package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	chat "google.golang.org/api/chat/v1"
)

//Grabbing config information for contact with Hangouts Chat's api.
var serviceKeyPath = os.Getenv("JIRABOT_SVC_KEY_PATH")

// GetAuthdChatClient consumes the service key given by google to perform
// an authorization.
func GetAuthdChatClient() *chat.Service {
	ctx := context.Background()

	data, err := ioutil.ReadFile(serviceKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := google.CredentialsFromJSON(
		ctx,
		data,
		"https://www.googleapis.com/auth/chat.bot",
	)
	if err != nil {
		log.Fatal(err)
	}

	service, err := chat.New(oauth2.NewClient(ctx, creds.TokenSource))
	if err != nil {
		log.Fatal("Failed to create chat service: " + err.Error())
	}

	return service
}
