package ipc

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/royiro10/cogo/models"
)

func SendMsg(conn net.Conn, msg models.CogoMessage) error {
	msgRaw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(conn, string(msgRaw)); err != nil {
		return err
	}

	return nil
}
