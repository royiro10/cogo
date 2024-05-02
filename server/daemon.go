package server

import (
	"context"
	"errors"
	"net"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/ipc"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/services"
	"github.com/royiro10/cogo/util"
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

	util.RedirectStderr(logger.LogFile)

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

	msg, err := ipc.ReciveMsg(conn)
	if err != nil {
		logger.Error("Error reading message", "err", err)
		return
	}

	logger.Info("Background process received message", "request", msg)

	switch req := msg.(type) {
	case *models.ExecuteRequest:
		daemon.commandService.HandleCommand(req)
		daemon.Ack(conn)
		return
	case *models.KillRequest:
		daemon.commandService.HandleKill(req)
		daemon.Ack(conn)
		return
	default:
		errMsg := "unkown request type"
		logger.Error(errMsg, "request", msg, "type", req)
		daemon.Err(conn, errors.New(errMsg))
	}
}

func (daemon *CogoDaemon) Ack(conn net.Conn) {
	daemon.sendResponse(conn, models.NewAckResponse())
}

func (daemon *CogoDaemon) Err(conn net.Conn, err error) {
	daemon.sendResponse(conn, models.NewErrResponse(err))
}

func (daemon *CogoDaemon) sendResponse(conn net.Conn, msg models.CogoMessage) {
	if err := ipc.SendMsg(conn, msg); err != nil {
		daemon.logger.Error("failed to send", "msg", msg)
	}
}
