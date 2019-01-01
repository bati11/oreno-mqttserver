package packet

import (
	"bufio"
	"encoding/binary"
	"io"
	"regexp"

	"github.com/pkg/errors"
)

type ConnectPayload struct {
	ClientID string
}

var clientIDRegex = regexp.MustCompile("^[a-zA-Z0-9-|]*$")

func ToConnectPayload(r *bufio.Reader) (ConnectPayload, error) {
	lengthBytes := make([]byte, 2)
	_, err := io.ReadFull(r, lengthBytes)
	if err != nil {
		return ConnectPayload{}, err
	}
	length := binary.BigEndian.Uint16(lengthBytes)

	clientIDBytes := make([]byte, length)
	_, err = io.ReadFull(r, clientIDBytes)
	if err != nil {
		return ConnectPayload{}, err
	}
	clientID := string(clientIDBytes)
	if len(clientID) < 1 || len(clientID) > 23 {
		return ConnectPayload{}, errors.New("ClientID length is invalid")
	}
	if !clientIDRegex.MatchString(clientID) {
		return ConnectPayload{}, errors.New("clientId format shoud be \"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ\"")
	}
	return ConnectPayload{ClientID: clientID}, nil
}
