package services

import (
	"bufio"
	"context"
	"encoding/json"
	"net"

	"github.com/royiro10/cogo/util"
)

type Daemon interface {
	Start(ctx context.Context)
}

type CogoDaemon struct {
	logger         *util.Logger
	commandService *CommandService
}

func CreateCogoDaemon(logger *util.Logger, commandService *CommandService) Daemon {
	d := &CogoDaemon{
		logger:         logger,
		commandService: commandService,
	}

	return d
}

func (daemon *CogoDaemon) Start(ctx context.Context) {
	logger := daemon.logger

	logger.Info("Daemon is running...")
	listener, closeFunc := util.MakeListener(logger)
	defer closeFunc()

	if listener == nil {
		logger.Error("can not start listening to message")
		return
	}

	logger.Debug("started socket server for IPC")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			daemon.handleMessage(listener)
		}
	}
}

func (daemon *CogoDaemon) handleMessage(listener net.Listener) {
	logger := daemon.logger

	logger.Debug("accepting connections")
	conn, err := listener.Accept()
	if err != nil {
		logger.Error("Failed to accept connection", "err", err)
		return
	}
	defer conn.Close()

	rawMessage, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		logger.Error("Error reading message", "err", err)
		return
	}

	logger.Info("Background process received message", "message", rawMessage)

	var message CommandParameters
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		logger.Error("Error parsing message", "err", err)
		return
	}

	daemon.commandService.HandleCommand(&message)
}
