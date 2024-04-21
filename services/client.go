package services

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/royiro10/cogo/util"
)

type CogoClient struct {
	logger      *util.Logger
	ipcClient   net.Conn
	releaseFunc util.IDisposable
}

func CreateCogoClient(logger *util.Logger) *CogoClient {
	ipcClient, releaseFunc := util.MakeIpcClient(logger)
	if ipcClient == nil {
		releaseFunc()
		util.LogErrorFatal(logger, "could not create connection to COGO", nil)
	}

	cogoClient := &CogoClient{
		ipcClient:   ipcClient,
		logger:      logger,
		releaseFunc: releaseFunc,
	}

	return cogoClient
}

func (client *CogoClient) Run(cp *CommandParameters) {
	msg, err := json.Marshal(cp)
	if err != nil {
		client.logger.Error("failed to serilize command request", "err", err)
		return
	}

	fmt.Fprintln(client.ipcClient, string(msg))
	client.logger.Debug("send command", "command", msg)
}
