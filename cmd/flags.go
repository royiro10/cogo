package cmd

import (
	"flag"

	"github.com/royiro10/cogo/models"
)

func InitFlags() *models.CogoCLIFlags {
	cogoFlags := &models.CogoCLIFlags{}

	sessionMsg := "Sppcify a session to interact with"
	flag.StringVar(&cogoFlags.Session, "session", "", sessionMsg)
	flag.StringVar(&cogoFlags.Session, "s", "", sessionMsg)

	isLogging := "Should log"
	flag.BoolVar(&cogoFlags.IsLogging, "logger", false, isLogging)
	flag.BoolVar(&cogoFlags.IsLogging, "l", false, isLogging)

	isStreamMsg := "Should follow and return print output in stream"
	flag.BoolVar(&cogoFlags.IsStream, "follow", false, isStreamMsg)
	flag.BoolVar(&cogoFlags.IsStream, "f", false, isStreamMsg)

	flag.Parse()

	return cogoFlags
}
