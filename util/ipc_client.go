package util

import (
	"net"
)

type IpcClient struct {
	Conn        net.Conn
	ReleaseFunc IDisposable
}
