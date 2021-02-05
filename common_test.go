package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTicketID(t *testing.T) {
	BotName = "@ABCBOT"
	wantedID := "XXX-1234"

	testCases := []struct {
		name, wantedID, text string
	}{
		{
			name:     "Find ID if called at the beginning",
			wantedID: wantedID,
			text:     fmt.Sprintf("%s %s", BotName, wantedID),
		},
		{
			name:     "Find ID if called at the end",
			wantedID: wantedID,
			text:     fmt.Sprintf("So what do you think we should size %s %s", BotName, wantedID),
		},
		{
			name:     "Find ID if called in the middle",
			wantedID: wantedID,
			text:     fmt.Sprintf("Use it like this %s %s, calling the bot then giving it a ticket", BotName, wantedID),
		},
		{
			name:     "Return only first ID",
			wantedID: wantedID,
			text:     fmt.Sprintf("Use it like this %s %s %s, calling the bot then giving it a ticket", BotName, wantedID, "BAD-420"),
		},
		{
			name:     "Return black if only called bot",
			wantedID: "",
			text:     fmt.Sprintf("%s", BotName),
		},
		{
			name:     "Return blank if no ID present in text",
			wantedID: "",
			text:     fmt.Sprintf("Use it like this %s calling the bot then giving it a ticket", BotName),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotID := getTicketID(ChatPayload{
				Message: Message{
					Text: tc.text,
				},
			})

			assert.Equal(t, tc.wantedID, gotID)
		})
	}
}
