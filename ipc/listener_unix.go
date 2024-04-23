//go:build linux || freebsd || darwin

package ipc

import (
	"net"
	"os"

	"github.com/royiro10/cogo/common"
)

func MakeIpcServer(logger *common.Logger) (*IpcServer, error) {
	listener, err := net.Listen("unix", COGO_CONN_UNIX)
	if err != nil {
		logger.Error("Failed to create Unix domain socket listener:", err)
		os.Remove(COGO_CONN_UNIX)
		return nil, err
	}

	server := &IpcServer{
		Listener: listener,
		ReleaseFunc: func() {
			os.Remove(COGO_CONN_UNIX)
			listener.Close()
		},
	}

	return server, nil
}
