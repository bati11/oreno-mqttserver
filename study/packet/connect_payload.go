package packet

import (
	"encoding/binary"
	"regexp"

	"github.com/pkg/errors"
)

type ConnectPayload struct {
	ClientID string
}

var clientIDRegex = regexp.MustCompile("^[a-zA-Z0-9-|]*$")

func ToConnectPayload(bs []byte) (ConnectPayload, error) {
	if len(bs) < 3 {
		return ConnectPayload{}, errors.New("payload length is invalid")
	}
	length := binary.BigEndian.Uint16(bs[0:2])
	var clientID string
	if len(bs) < 2+int(length) {
		return ConnectPayload{}, errors.New("specified length is not equals ClientID length")
	} else {
		clientID = string(bs[2 : 2+length])
	}
	if len(clientID) < 1 || len(clientID) > 23 {
		return ConnectPayload{}, errors.New("ClientID length is invalid")
	}
	if !clientIDRegex.MatchString(clientID) {
		return ConnectPayload{}, errors.New("clientId format shoud be \"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ\"")
	}
	return ConnectPayload{ClientID: clientID}, nil
}
