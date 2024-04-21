//go:build windows

package util

import (
	"net"
)

func MakeIpcClient(logger *Logger) (net.Conn, IDisposable) {
	conn, err := net.Dial("tcp", COGO_CONN_WIN32)
	if err != nil {
		logger.Error("Failed to connect to background process", "err", err, "addr", COGO_CONN_WIN32)
		return nil, func() {}
	}

	return conn, func() {
		conn.Close()
	}
}
