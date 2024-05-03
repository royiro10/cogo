package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/royiro10/cogo/models"
)

func ReciveMsg(conn net.Conn) (models.CogoMessage, error) {
	rawMessage, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading(ReadingBytes): %v", err)
	}

	var data models.BaseCogoMessage
	if err := json.Unmarshal(rawMessage, &data); err != nil {
		return nil, fmt.Errorf("error reading(UnmarshelBaseMsg): %v", err)
	}

	msg := models.GetMessage(*data.GetDetails())
	if msg == nil {
		return nil, fmt.Errorf("error reading(Unrecognized msg) %v", data.GetDetails())
	}

	if err := json.Unmarshal(rawMessage, msg); err != nil {
		return nil, fmt.Errorf("error reading(UnmarshelMsg::%s): %v", msg.GetDetails().Type, err)
	}

	return msg, nil
}
