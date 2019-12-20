package packet

import (
	"encoding/binary"
	"io"
	"regexp"
)

type ConnectPayload struct {
	ClientID string
}

var clientIDRegex = regexp.MustCompile("^[a-zA-Z0-9-|]*$")

func (reader *MQTTReader) readConnectPayload() (*ConnectPayload, error) {
	lengthBytes := make([]byte, 2)
	_, err := io.ReadFull(reader, lengthBytes)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(lengthBytes)

	clientIDBytes := make([]byte, length)
	_, err = io.ReadFull(reader, clientIDBytes)
	if err != nil {
		return nil, err
	}
	clientID := string(clientIDBytes)
	if len(clientID) < 1 || len(clientID) > 23 {
		return nil, RefusedByIdentifierRejected("ClientID length is invalid")
	}
	if !clientIDRegex.MatchString(clientID) {
		return nil, RefusedByIdentifierRejected("ClientId format shoud be \"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ\"")
	}
	return &ConnectPayload{ClientID: clientID}, nil
}
