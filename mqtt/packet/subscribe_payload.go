package packet

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

type SubscribePayload struct {
	TopicFilter string
	QoS         QoS
}

func ToSubscribePayload(r io.Reader) (SubscribePayload, error) {
	lengthBuf := make([]byte, 2)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return SubscribePayload{}, errors.New("")
	}
	length := binary.BigEndian.Uint16(lengthBuf)

	buf := make([]byte, length+1)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return SubscribePayload{}, errors.New("")
	}
	topicFilter := buf[:length]
	qos := buf[length]
	qosBit1 := refbit(qos, 1) > 0
	qosBit2 := refbit(qos, 0) > 0
	result := SubscribePayload{
		TopicFilter: string(topicFilter),
		QoS:         toQoS(qosBit1, qosBit2),
	}
	return result, nil
}
