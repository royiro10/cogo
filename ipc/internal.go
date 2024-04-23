package ipc

import (
	"net"

	"github.com/royiro10/cogo/common"
)

const COGO_CONN_WIN32 = "localhost:3001"
const COGO_CONN_UINX = "./cogo.sock"

type IpcClient struct {
	Conn        net.Conn
	ReleaseFunc common.IDisposable
}

type IpcServer struct {
	Listener    net.Listener
	ReleaseFunc common.IDisposable
}
