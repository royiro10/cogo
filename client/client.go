package client

import (
	"encoding/json"
	"fmt"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/ipc"
	"github.com/royiro10/cogo/services"
)

type CogoClient struct {
	logger    *common.Logger
	ipcClient *ipc.IpcClient
}

func CreateCogoClient(logger *common.Logger) *CogoClient {
	makeIpcClient := common.MakeRetryable[*ipc.IpcClient](func() (*ipc.IpcClient, error) {
		return ipc.MakeIpcClient(logger)
	}, 3)

	ipcClient, err := makeIpcClient()
	if err != nil {
		logger.Fatal(common.WrapedError{Msg: "could not create connection to COGO", Err: err})
	}

	cogoClient := &CogoClient{
		ipcClient: ipcClient,
		logger:    logger,
	}

	return cogoClient
}

func (client *CogoClient) Run(cp *services.CommandParameters) {
	msg, err := json.Marshal(cp)
	if err != nil {
		client.logger.Error("failed to serilize command request", "err", err)
		return
	}

	fmt.Fprintln(client.ipcClient.Conn, string(msg))
	client.logger.Debug("send command", "command", msg)
}

func (client *CogoClient) Close() {
	client.ipcClient.ReleaseFunc()
}
