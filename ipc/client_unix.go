//go:build linux || freebsd || darwin

package ipc

import (
	"net"

	"github.com/royiro10/cogo/common"
)

func MakeIpcClient(logger *common.Logger) (*IpcClient, error) {
	conn, err := net.Dial("unix", GetUnixConnection())
	if err != nil {
		logger.Error(
			"failed to connect to background process",
			"err",
			err,
			"addr",
			GetUnixConnection(),
		)
		return nil, err
	}

	return &IpcClient{
		Conn:        conn,
		ReleaseFunc: func() { conn.Close() },
	}, nil
}
