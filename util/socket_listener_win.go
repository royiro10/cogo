//go:build windows

package util

import (
	"net"
)

const COGO_CONN_WIN32 = "localhost:3001"

func MakeIpcListener(logger *Logger) (net.Listener, IDisposable) {
	listener, err := net.Listen("tcp", COGO_CONN_WIN32)
	if err != nil {
		logger.Error("Failed to start TCP server:", err)
		return nil, func() {}
	}

	closeHandler := func() {
		listener.Close()
	}

	return listener, closeHandler
}
