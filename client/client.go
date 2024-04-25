package client

import (
	"encoding/json"
	"fmt"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/ipc"
	"github.com/royiro10/cogo/models"
)

type CogoClient struct {
	logger    *common.Logger
	ipcClient *ipc.IpcClient
}

func CreateCogoClient(logger *common.Logger) *CogoClient {
	makeIpcClient := common.MakeRetryable(func() (*ipc.IpcClient, error) {
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

func (client *CogoClient) sendData(request interface{}) {
	msg, err := json.Marshal(request)
	if err != nil {
		client.logger.Error("failed to serilize request", "err", err)
		return
	}

	fmt.Fprintln(client.ipcClient.Conn, string(msg))
}

func (client *CogoClient) Run(request *models.ExecuteRequest) {
	client.logger.Debug("send run command", "request", request)
	client.sendData(request)
}

func (client *CogoClient) Kill(request *models.KillRequest) {
	client.logger.Debug("send kill command", "request", request)
	client.sendData(request)
}

func (client *CogoClient) Close() {
	client.ipcClient.ReleaseFunc()
}
