package main

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	chat "google.golang.org/api/chat/v1"
)

// ChatThreadKey is a arbitrary value to maintain thread
const ChatThreadKey = "JIRABOT_THREADKEY"

// ChatLogRoomID is the room the bot's logs will go to
var ChatLogRoomID = os.Getenv("CHAT_LOG_ROOM_ID")

// LogLevel denotes the specific level of logs to
// output
var LogLevel = os.Getenv("JIRABOT_LOG_LEVEL")

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

	_, err := chat.NewSpacesMessagesService(clw.Service).
		Create("spaces/"+ChatLogRoomID, msg).
		ThreadKey(ChatThreadKey).
		Do()
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
		level = logrus.WarnLevel
	}

	logger := logrus.New()
	logger.SetLevel(level)

	if level > logrus.WarnLevel {
		logger.SetReportCaller(true)
	}

	logger.SetOutput(io.MultiWriter(writers...))

	return logger
}
