package main

import (
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

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

	// Removes any extra characters from the end of the ticket
	re := regexp.MustCompile(`\D$`)
	for re.MatchString(cp.Message.Args) {
		cp.Message.Args = string(cp.Message.Args[:len(cp.Message.Args)-1])
	}

	spew.Dump(cp.Message.Args)
}
