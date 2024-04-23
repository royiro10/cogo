//go:build windows

package ipc

import (
	"net"

	"github.com/royiro10/cogo/common"
)

func MakeIpcServer(logger *common.Logger) (*IpcServer, error) {
	listener, err := net.Listen("tcp", COGO_CONN_WIN32)
	if err != nil {
		logger.Error("Failed to start TCP server:", err)
		return nil, err
	}

	server := &IpcServer{
		Listener: listener,
		ReleaseFunc: func() {
			listener.Close()
		},
	}

	return server, nil
}
