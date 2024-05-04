package cmd

import (
	"flag"

	"github.com/royiro10/cogo/models"
)

func InitFlags() *models.CogoCLIFlags {
	cogoFlags := &models.CogoCLIFlags{}

	sessionMsg := "sepcify a session to interact with"
	flag.StringVar(&cogoFlags.Session, "session", "", sessionMsg)
	flag.StringVar(&cogoFlags.Session, "s", "", sessionMsg)

	isLogging := "shold log"
	flag.BoolVar(&cogoFlags.IsLogging, "logger", false, isLogging)
	flag.BoolVar(&cogoFlags.IsLogging, "l", false, isLogging)

	isStreamMsg := "should follow and return print output in stream"
	flag.BoolVar(&cogoFlags.IsStream, "follow", false, isStreamMsg)
	flag.BoolVar(&cogoFlags.IsStream, "f", false, isStreamMsg)

	flag.Parse()

	return cogoFlags
}
