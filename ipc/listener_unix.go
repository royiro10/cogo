//go:build linux || freebsd || darwin

package ipc

import (
	"net"
	"os"

	"github.com/royiro10/cogo/common"
)

func MakeIpcServer(logger *common.Logger) (*IpcServer, error) {
	listener, err := net.Listen("unix", GetUnixConnection())
	if err != nil {
		logger.Error("failed to create Unix domain socket listener:", err)
		os.Remove(GetUnixConnection())
		return nil, err
	}

	server := &IpcServer{
		Listener: listener,
		ReleaseFunc: func() {
			os.Remove(GetUnixConnection())
			listener.Close()
		},
	}

	return server, nil
}
