package main

import (
	"os"
	"regexp"
	"strings"
)

// BotName is the name of the bot
var BotName = os.Getenv("JIRABOT_BOT_NAME")

// ChatPayload consumes the essential information from the data sent by
// google chat
type ChatPayload struct {
	Type    string  `json:"type"`
	Message Message `json:"message"`
}

// Message holds the text passed to the bot
type Message struct {
	Text string `json:"text"`
}

// GetTicketID Reads the payload and returns the
func getTicketID(payload ChatPayload) string {
	botIdx := strings.Index(payload.Message.Text, BotName)
	args := strings.Fields(payload.Message.Text[botIdx:])

	ticketID := ""
	if len(args) > 1 {
		ticketID = args[1]
	}

	return regexp.MustCompile(`[a-zA-Z]+-\d+`).FindString(ticketID)
}
