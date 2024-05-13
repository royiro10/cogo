package client

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/ipc"
	"github.com/royiro10/cogo/models"
)

type CogoClient struct {
	logger    *common.Logger
	ipcClient *ipc.IpcClient
}

type CogoClientChan chan string

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

func (client *CogoClient) Status() error {
	client.sendData(models.NewListSessionsRequest())

	var sessions []string
	var waitingForAck bool = true

	for waitingForAck {
		msg, err := ipc.ReciveMsg(client.ipcClient.Conn)
		if err != nil {
			if err == ipc.ErrConnectionClosed {
				client.logger.Debug("connection closed from the daemon, assume stream has ended")
				return nil
			}

			return err
		}

		switch response := msg.(type) {
		case *models.ListSessionsResponse:
			client.logger.Info("recive sessions ids list", "listCount", len(response.Sessions))
			sessions = response.Sessions

		case *models.AckResponse:
			client.logger.Debug("recive ack response to end listening for outputs")
			waitingForAck = false

		default:
			client.logger.Error("unkown response type", "response", response)
			return fmt.Errorf("unkown response type: %s, %v", response.GetDetails().Type, response)
		}
	}

	ipcClient, err := ipc.MakeIpcClient(client.logger)
	if err != nil {
		return fmt.Errorf("FUCK")
	}

	client.ipcClient = ipcClient
	errChan := make(chan error)
	var wg sync.WaitGroup

	wg.Add(len(sessions))

	for _, s := range sessions {
		client.logger.Info("status for session", "session", s)
		go func(session string) {
			defer wg.Done()

			ipcClient, err := ipc.MakeIpcClient(client.logger)
			defer ipcClient.ReleaseFunc()

			if err != nil {
				errChan <- err
				return
			}

			// TODO: change this to come from request param
			req := models.NewStatusRequest(session, false)
			client.logger.Debug("send command", "type", req.GetDetails().Type, "request", req)

			if err := ipc.SendMsg(ipcClient.Conn, req); err != nil {
				client.logger.Error("could not send request", "reason", err.Error())
			}

			for {
				client.logger.Debug("waiting for message")
				msg, err := ipc.ReciveMsg(ipcClient.Conn)

				if err != nil {
					if err == ipc.ErrConnectionClosed {
						client.logger.Debug("connection closed from the daemon, assume stream has ended")
						return
					}

					client.logger.Error(err.Error())
					errChan <- err
				}

				client.logger.Debug("recived message", "msg", msg.GetDetails())

				switch response := msg.(type) {
				case *models.StatusResponse:
					client.logger.Info("recived status", "sessionId", response.SessionId, "status", response.Status)
					statusTimeFormat := "2006-01-02 15:04:05"

					baseFormat := "{{.LastActionTime}} {{.SessionId}}: {{.SessionStatus}} => {{.LastCommand}}"
					verboseFormat := "executed:{{.ExecutedCommandsCount}} queue:{{.CommandsToExecuteQueueSize}} stdout:{{.OutputViewSize}}"

					format := fmt.Sprintf("%s\n\t(%s)", baseFormat, verboseFormat)

					// _ = "{{.LastActionTime}} {{.SessionId}}: {{.SessionStatus}} | executed:{{.ExecutedCommandsCount}} queue:{{.CommandsToExecuteQueueSize}} stdout:{{.OutputViewSize}}\n\t{{.LastCommand}}"

					sessionStatuMsg := formatMsg(format, map[string]interface{}{
						"LastActionTime":             response.Status.LastActionTime.Format(statusTimeFormat),
						"SessionId":                  response.SessionId,
						"SessionStatus":              response.Status.SessionStatus,
						"ExecutedCommandsCount":      response.Status.ExecutedCommandCount,
						"CommandsToExecuteQueueSize": response.Status.ExecuteQueueSize,
						"OutputViewSize":             response.Status.OutputViewSize,
						"LastCommand":                response.Status.LastCommand,
					})

					fmt.Println(sessionStatuMsg)

				case *models.AckResponse:
					client.logger.Debug("recive ack response to end listening for outputs")
					return

				default:
					client.logger.Error("unkown response type", "response", response)
					errChan <- fmt.Errorf("unkown response type: %s, %v", response.GetDetails().Type, response)
				}
			}
		}(s)
	}

	ctx, cancel := context.WithCancel(context.TODO())

	go func() {
		for {
			select {
			case err := <-errChan:
				client.logger.Error(err.Error())
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
	cancel()
	return nil
}

func (client *CogoClient) ListSessions(request *models.ListSessionsRequest) error {
	client.sendData(request)

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
		case *models.ListSessionsResponse:
			delim := ", "
			fmt.Printf("sessions list:\n\t%s\n", strings.Join(response.Sessions, delim))

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

func formatMsg(fmt string, args map[string]interface{}) (str string) {
	var msg bytes.Buffer

	tmpl, err := template.New("formmatedMsg").Parse(fmt)

	if err != nil {
		return fmt
	}

	tmpl.Execute(&msg, args)
	return msg.String()
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
