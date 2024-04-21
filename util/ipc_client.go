//go:build linux || freebsd || darwin

package util

import (
	"net"
)

func MakeIpcClient(logger *util.Logger) (net.Conn, IDisposable) {
	conn, err := net.Dial("unix", COGO_CONN_UINX)
	if err != nil {
		logger.Error("Failed to connect to background process", "err", err, "addr", COGO_CONN_UINX)
		return nil, func() {}
	}

	return conn, func() {
		conn.Close()
	}
}
