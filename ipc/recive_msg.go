package ipc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/royiro10/cogo/models"
)

var ErrConnectionClosed error = errors.New("connection closed")

func ReciveMsg(conn net.Conn) (models.CogoMessage, error) {
	ipcPacket, err := reciveIpcPacket(conn)
	if err != nil {
		return nil, err
	}

	msg := models.GetMessage(*ipcPacket.Message.GetDetails())
	if msg == nil {
		return nil, fmt.Errorf("error reading(Unrecognized msg) %v", ipcPacket.Message.GetDetails())
	}

	if err := json.Unmarshal(ipcPacket.Raw, msg); err != nil {
		return nil, fmt.Errorf("error reading(UnmarshelMsg::%s): %v", msg.GetDetails().Type, err)
	}

	return msg, nil
}

func reciveIpcPacket(conn net.Conn) (*IpcPacket, error) {
	headerSize := binary.Size(ipcHeaderDefinition{})
	headerBytes := make([]byte, headerSize)
	if _, err := io.ReadFull(conn, headerBytes); err != nil {
		return nil, err
	}

	packet := IpcPacket{
		Header:  ipcHeaderDefinition{},
		Raw:     nil,
		Message: &models.BaseCogoMessage{},
	}

	if err := binary.Read(bytes.NewReader(headerBytes), binary.BigEndian, &packet.Header); err != nil {
		return &packet, err
	}

	packet.Raw = make([]byte, packet.Header.MessageSize)
	if _, err := io.ReadFull(conn, packet.Raw); err != nil {
		return &packet, err
	}

	if err := json.Unmarshal(packet.Raw, packet.Message); err != nil {
		return &packet, fmt.Errorf("error reading(UnmarshelBaseMsg): %v", err)
	}

	return &packet, nil
}
