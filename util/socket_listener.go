//go:build linux || freebsd || darwin

package util

import (
	"net"
	"os"
)

const COGO_CONN_UINX = "./cogo.sock"

func MakeListener(logger *util.Logger) (net.Listener, IDisposable) {
	listener, err := net.Listen("unix", COGO_CONN_UINX)
	if err != nil {
		logger.Error("Failed to create Unix domain socket listener:", err)
		return nil, func() {
			os.Remove(COGO_CONN_UINX)
		}
	}

	closeHandler := func() {
		os.Remove(COGO_CONN_UINX)
		listener.Close()
	}

	return listener, closeHandler
}
