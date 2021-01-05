package main

import "os"

type botConfig struct {
	certFile, certKeyFile string
	useSSL                bool
}

func loadConfig() (config botConfig) {
	config = botConfig{
		certFile:    os.Getenv("JIRABOT_CERT_FILE"),
		certKeyFile: os.Getenv("JIRABOT_CERT_KEY_FILE"),
		useSSL:      false,
	}

	if os.Getenv("JIRABOT_USE_SSL") == "true" {
		config.useSSL = true
	}

	return
}
