package ipc

import (
	"net"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
)

const COGO_CONN_WIN32 = "localhost:3001"
const COGO_CONN_UNIX = "./cogo.sock"

type IpcClient struct {
	Conn        net.Conn
	ReleaseFunc common.IDisposable
}

type IpcServer struct {
	Listener    net.Listener
	ReleaseFunc common.IDisposable
}

var IPCPacketVersion = 1

type ipcHeaderDefinition struct {
	Version     uint16
	MessageSize uint64 // not incuding header
}

type IpcPacket struct {
	Header  ipcHeaderDefinition
	Raw     []byte
	Message *models.BaseCogoMessage
}
