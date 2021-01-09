package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	chat "google.golang.org/api/chat/v1"
)

// LogLevel denotes the specific level of logs to
// output
var LogLevel = os.Getenv("JIRABOT_LOG_LEVEL")

// ChatParentKey is a arbitrary value to maintain thread
const ChatParentKey = "JIRABOT_THREADKEY"

// ChatLogRoomID is the room the bot's logs will go to
var ChatLogRoomID = os.Getenv("CHAT_LOG_ROOM_ID")

// ChatLogWriter is a wrapper for a chat.Service, it's purpose is to be an
// io.Writer to send logs to a chatroom.
type ChatLogWriter struct {
	*chat.Service
}

// Authorize authorizes the logger to be allow log transmission
// to gchat
func (clw *ChatLogWriter) Authorize() {
	clw.Service = GetAuthdChatClient()
}

func (clw ChatLogWriter) Write(data []byte) (int, error) {
	msg := &chat.Message{
		Text: string(data),
	}

	msgSvc := chat.NewSpacesMessagesService(clw.Service)
	msgCall := msgSvc.Create(ChatParentKey, msg)

	_, err := msgCall.Do()
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

//StartLogger function initializes the logger
//given the specified configurations
func StartLogger(writers ...io.Writer) *logrus.Logger {
	level, e := logrus.ParseLevel(LogLevel)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetReportCaller(true)

	if len(writers) > 0 {
		mw := io.MultiWriter(writers...)
		logger.SetOutput(mw)
	} else {
		logger.SetOutput(os.Stderr)
	}

	return logger
}
