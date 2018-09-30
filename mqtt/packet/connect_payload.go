package packet

import (
	"encoding/binary"
	"errors"
	"regexp"
)

type ConnectPayload struct {
	ClientID string
}

var clientIDRegex = regexp.MustCompile("^[a-zA-Z0-9]*$")

var (
	ErrConnectPayloadLength = errors.New("payload length of connect packet is invalid")
	ErrClientIDFormat       = errors.New("clientId format shoud be \"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ\"")
)

func ToConnectPayload(bs []byte) (ConnectPayload, error) {
	if len(bs) < 3 {
		return ConnectPayload{}, ErrConnectPayloadLength
	}
	length := binary.BigEndian.Uint16(bs[0:2])
	var s string
	if len(bs) < 2+int(length) {
		s = string(bs[2:])
	} else {
		s = string(bs[2 : 2+length])
	}
	if len(s) < 1 || len(s) > 23 {
		return ConnectPayload{}, ErrClientIDFormat
	}
	if !clientIDRegex.MatchString(s) {
		return ConnectPayload{}, ErrClientIDFormat
	}
	result := ConnectPayload{
		ClientID: s,
	}
	return result, nil
}
