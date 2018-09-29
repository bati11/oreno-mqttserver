package packet

import (
	"encoding/binary"
)

type ConnectPayload struct {
	ClientID string
}

func ToConnectPayload(bs []byte) ConnectPayload {
	length := binary.BigEndian.Uint16(bs[0:2])
	return ConnectPayload{
		ClientID: string(bs[2 : 2+length]),
	}
}
