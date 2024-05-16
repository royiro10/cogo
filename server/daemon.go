package server

import (
	"context"
	"errors"
	"net"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/ipc"
	"github.com/royiro10/cogo/models"
	services "github.com/royiro10/cogo/services/commands"
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
		_ = daemon.Ack(conn)
		return
	case *models.KillRequest:
		daemon.commandService.HandleKill(req)
		_ = daemon.Ack(conn)
		return
	case *models.OutputRequest:
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		for output := range daemon.commandService.HandleOutput(req, ctx) {
			daemon.logger.Info("Sending", "output", output)

			if err := daemon.Output(conn, output); err != nil {
				daemon.logger.Warn("a connection has been closed while streaming, stop streaming")
				cancel()
			}
		}

		_ = daemon.Ack(conn)
	default:
		errMsg := "unknown request type"
		logger.Error(errMsg, "request", msg, "type", req)
		if err := daemon.Err(conn, errors.New(errMsg)); err != nil {
			logger.Error("failed to send error response to connection", "remoteAddr", conn.RemoteAddr().String())
		}
	}
}

func (daemon *CogoDaemon) Ack(conn net.Conn) error {
	return daemon.sendResponse(conn, models.NewAckResponse())
}

func (daemon *CogoDaemon) Err(conn net.Conn, err error) error {
	return daemon.sendResponse(conn, models.NewErrResponse(err))
}

func (daemon *CogoDaemon) Output(conn net.Conn, output *models.StdLine) error {
	return daemon.sendResponse(conn, models.NewOutputResponse(output))
}

func (daemon *CogoDaemon) sendResponse(conn net.Conn, msg models.CogoMessage) error {
	if err := ipc.SendMsg(conn, msg); err != nil {
		daemon.logger.Error("failed to send", "msg", msg)
		return err
	}

	return nil
}
