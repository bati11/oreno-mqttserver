package packet

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

type TopicFilterPair struct {
	Filter string
	QoS    uint8
}
type SubscribePayload struct {
	TopicFilterPairs []*TopicFilterPair
}

func ToSubscribePayload(fixedHeader FixedHeader, variableHeader SubscribeVariableHeader, r *bufio.Reader) (SubscribePayload, error) {
	if fixedHeader.PacketType != SUBSCRIBE {
		return SubscribePayload{}, fmt.Errorf("packet type is invalid. it got is %v", fixedHeader.PacketType)
	}

	var topicFilterPairs []*TopicFilterPair

	remain := fixedHeader.RemainingLength - variableHeader.Length()
	for remain > 0 {
		length, err := extractLength(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			return SubscribePayload{}, err
		}
		remain = remain - 2 - uint(length)

		if remain < 0 {
			break
		}

		bs := make([]byte, length)
		_, err = io.ReadFull(r, bs)
		if err != nil {
			return SubscribePayload{}, err
		}
		topicFilter := string(bs)

		qos, err := extractQoS(r)
		if err != nil {
			return SubscribePayload{}, err
		}
		remain--

		topicFilterPairs = append(topicFilterPairs, &TopicFilterPair{topicFilter, qos})
	}
	return SubscribePayload{topicFilterPairs}, nil
}

func extractLength(r *bufio.Reader) (uint16, error) {
	lengthMSB, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	lengthLSB, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	length := binary.BigEndian.Uint16([]byte{lengthMSB, lengthLSB})
	return length, nil
}

func extractQoS(r *bufio.Reader) (uint8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	qos1 := refbit(b, 1)
	qos2 := refbit(b, 0)
	qos := qos2 << qos1
	return qos, nil
}
