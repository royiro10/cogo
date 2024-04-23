package server

import (
	"bufio"
	"context"
	"encoding/json"
	"net"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/ipc"
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

	var message services.CommandParameters
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		logger.Error("Error parsing message", "err", err)
		return
	}

	logger.Info("Background process received message", "message", message)

	daemon.commandService.HandleCommand(&message)
}
