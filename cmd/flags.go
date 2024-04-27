package cmd

import (
	"flag"

	"github.com/royiro10/cogo/models"
)

func InitFlags() *models.CogoCLIFlags {
	cogoFlags := &models.CogoCLIFlags{}

	flag.StringVar(&cogoFlags.Session, "session", "", "sepcify a session to interact with")
	flag.StringVar(&cogoFlags.Session, "s", "", "sepcify a session to interact with")

	flag.BoolVar(&cogoFlags.IsLogging, "logger", false, "shold log")
	flag.BoolVar(&cogoFlags.IsLogging, "l", false, "shold log")

	flag.Parse()

	return cogoFlags
}
