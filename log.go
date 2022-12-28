package main

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/global"
)

// Create a logrus logger.
// Outputs to both stdout and log file
func MakeLogger(logFilePath string) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err.Error())
	}

	// Output to both stdout and log file
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	SetLogLevelFromEnv()
}

// Get & set logrus log level from environment variable
func SetLogLevelFromEnv() {
	switch global.LOG_LEVEL {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.PanicLevel)
	}
}
