package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// LogLevel denotes the specific level of logs to
// output
var LogLevel = os.Getenv("JIRABOT_LOG_LEVEL")

// LogFile is the location to log to. If empty
// log to stderr
var LogFile = os.Getenv("JIRABOT_LOG_FILE")

//StartLogger function initializes the logger
//given the specified configurations
func StartLogger() *logrus.Logger {
	level, e := logrus.ParseLevel(LogLevel)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	logFH := os.Stderr
	if LogFile != "" {
		fd, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			fmt.Printf("Coudln't open log file: %v\nLogging to STDERR.", err)
		} else {
			logFH = fd
		}
	}

	logger := logrus.New()
	logger.Out = logFH
	logger.Level = level
	logger.SetReportCaller(true)

	return logger
}
