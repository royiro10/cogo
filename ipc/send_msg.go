package ipc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"

	"github.com/royiro10/cogo/models"
)

func SendMsg(conn net.Conn, msg models.CogoMessage) error {
	rawMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	header := ipcHeaderDefinition{
		Version:     uint16(IPCPacketVersion),
		MessageSize: uint64(len(rawMsg)),
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, header); err != nil {
		return err
	}

	if _, err := buf.Write(rawMsg); err != nil {
		return err
	}

	_, err = conn.Write(buf.Bytes())
	return err
}
