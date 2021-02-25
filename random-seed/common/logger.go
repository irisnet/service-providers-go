package common

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Logger is a logger instance
var Logger = log.New()

func init() {
	Logger.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	Logger.SetOutput(os.Stdout)

	Logger.SetLevel(log.InfoLevel)
}
