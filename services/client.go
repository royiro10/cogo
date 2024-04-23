package services

import (
	"encoding/json"
	"fmt"

	"github.com/royiro10/cogo/util"
)

type CogoClient struct {
	logger    *util.Logger
	ipcClient *util.IpcClient
}

func CreateCogoClient(logger *util.Logger) *CogoClient {
	makeIpcClient := util.MakeRetryable[*util.IpcClient](func() (*util.IpcClient, error) {
		return util.MakeIpcClient(logger)
	}, 3)

	ipcClient, err := makeIpcClient()
	if err != nil {
		util.LogErrorFatal(logger, "could not create connection to COGO", err)
	}

	cogoClient := &CogoClient{
		ipcClient: ipcClient,
		logger:    logger,
	}

	return cogoClient
}

func (client *CogoClient) Run(cp *CommandParameters) {
	msg, err := json.Marshal(cp)
	if err != nil {
		client.logger.Error("failed to serilize command request", "err", err)
		return
	}

	fmt.Fprintln(client.ipcClient.Conn, string(msg))
	client.logger.Debug("send command", "command", msg)
}
