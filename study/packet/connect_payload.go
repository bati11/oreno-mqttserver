package packet

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

type ConnectPayload struct {
	ClientID string
}

func ToConnectPayload(bs []byte) (ConnectPayload, error) {
	if len(bs) < 3 {
		return ConnectPayload{}, errors.New("payload length is invalid")
	}
	length := binary.BigEndian.Uint16(bs[0:2])
	var clientID string
	if len(bs) < 2+int(length) {
		clientID = string(bs[2:])
	} else {
		clientID = string(bs[2 : 2+length])
	}
	if len(clientID) < 1 || len(clientID) > 23 {
		return ConnectPayload{}, errors.New("ClientID length is invalid")
	}
	return ConnectPayload{ClientID: clientID}, nil
}
