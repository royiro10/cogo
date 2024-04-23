//go:build linux || freebsd || darwin

package util

import (
	"net"
)

func MakeIpcClient(logger *Logger) (*IpcClient, error) {
	conn, err := net.Dial("unix", COGO_CONN_UINX)
	if err != nil {
		logger.Error("Failed to connect to background process", "err", err, "addr", COGO_CONN_UINX)
		return nil, err
	}

	return &IpcClient{
		Conn:        conn,
		ReleaseFunc: func() { conn.Close() },
	}, nil
}
