package server

import (
	"bufio"
	"context"
	"encoding/json"
	"net"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/ipc"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/services"
)

type Daemon interface {
	Start(ctx context.Context)
}

type CogoDaemon struct {
	logger         *common.Logger
	commandService *services.CommandService
}

func CreateCogoDaemon(logger *common.Logger, commandService *services.CommandService) Daemon {
	d := &CogoDaemon{
		logger:         logger,
		commandService: commandService,
	}

	return d
}

func (daemon *CogoDaemon) Start(ctx context.Context) {
	logger := daemon.logger

	logger.Info("Daemon is running...")
	server, err := ipc.MakeIpcServer(logger)
	if err != nil {
		logger.Error("can not start listening to message")
		return
	}
	defer server.ReleaseFunc()

	logger.Debug("started socket server for IPC", "addr", server.Listener.Addr().String())

	for {
		select {
		case <-ctx.Done():
			logger.Info("stop recived cancel from ctx")
			return
		default:
			logger.Debug("accepting connections")
			conn, err := server.Listener.Accept()
			if err != nil {
				daemon.logger.Error("Error accepting connection", "error", err)
			}

			go daemon.handleMessage(conn)
		}
	}
}

func (daemon *CogoDaemon) handleMessage(conn net.Conn) {
	logger := daemon.logger
	defer conn.Close()

	rawMessage, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		logger.Error("Error reading message", "err", err)
		return
	}

	var data struct {
		Details models.CogoMessageDetails
	}

	if err := json.Unmarshal(rawMessage, &data); err != nil {
		logger.Error("Error reading message", "err", err)
		return
	}

	logger.Info("hi", data.Details.Type, "data", rawMessage)
	request := models.GetRequest(data.Details.Type)
	if err := json.Unmarshal(rawMessage, request); err != nil {
		logger.Error("Error reading message", "err", err)
		return
	}

	logger.Info("Background process received message", "request", request)
	switch req := request.(type) {
	case *models.ExecuteRequest:
		daemon.commandService.HandleCommand(req)
		return
	case *models.KillRequest:
		daemon.commandService.HandleKill(req)
		return
	default:
		logger.Error("unkown request type", "request", request, "type", req)
	}
}
