package client

import (
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

func (client *CogoClient) Run(request *models.ExecuteRequest) error {
	client.sendData(request)
	return client.ensureAck()
}

func (client *CogoClient) Kill(request *models.KillRequest) error {
	client.sendData(request)
	return client.ensureAck()
}

func (client *CogoClient) Output(request *models.OutputRequest) error {
	client.sendData(request)

	fmt.Printf("output for session: %s\n", request.SessionId)

	for {
		msg, err := ipc.ReciveMsg(client.ipcClient.Conn)
		if err != nil {
			if err == ipc.ErrConnectionClosed {
				client.logger.Debug("connection closed from the daemon, assume stream has ended")
				return nil
			}

			return err
		}

		switch response := msg.(type) {
		case *models.OutputResponse:
			fmt.Println(response.StdLine.Data)

		case *models.AckResponse:
			client.logger.Debug("recive ack response to end listening for outputs")
			return nil

		default:
			client.logger.Error("unkown response type", "response", response)
			return fmt.Errorf("unkown response type: %s, %v", response.GetDetails().Type, response)
		}
	}
}

func (client *CogoClient) Close() {
	client.ipcClient.ReleaseFunc()
}

func (client *CogoClient) sendData(request models.CogoMessage) {
	client.logger.Debug("send command", "type", request.GetDetails().Type, "request", request)

	if err := ipc.SendMsg(client.ipcClient.Conn, request); err != nil {
		client.logger.Error("could not send request", "reason", err.Error())
	}
}

func (client *CogoClient) ensureAck() error {
	msg, err := ipc.ReciveMsg(client.ipcClient.Conn)
	if err != nil {
		return err
	}

	client.logger.Debug("recived msg", "msg", msg)

	if msg.GetDetails().Type != models.AckResponseDetails.Type {
		return fmt.Errorf("unexpected response %v", msg.GetDetails())
	}

	return nil
}
