//go:build linux || freebsd || darwin

package ipc

import (
	"net"

	"github.com/royiro10/cogo/common"
)

func MakeIpcClient(logger *common.Logger) (*IpcClient, error) {
	conn, err := net.Dial("unix", COGO_CONN_UNIX)
	if err != nil {
		logger.Error("Failed to connect to background process", "err", err, "addr", COGO_CONN_UNIX)
		return nil, err
	}

	return &IpcClient{
		Conn:        conn,
		ReleaseFunc: func() { conn.Close() },
	}, nil
}
