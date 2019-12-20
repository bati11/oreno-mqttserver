package packet

import (
	"encoding/binary"
	"io"
)

type TopicFilterPair struct {
	Filter string
	QoS    uint8
}
type SubscribePayload struct {
	TopicFilterPairs []*TopicFilterPair
}

func (reader *MQTTReader) readSubscribePayload(payloadLength uint) (*SubscribePayload, error) {
	var topicFilterPairs []*TopicFilterPair

	remain := payloadLength
	for remain > 0 {
		length, err := extractLength(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		remain = remain - 2 - uint(length)

		if remain < 0 {
			break
		}

		bs := make([]byte, length)
		_, err = io.ReadFull(reader, bs)
		if err != nil {
			return nil, err
		}
		topicFilter := string(bs)

		qos, err := extractQoS(reader)
		if err != nil {
			return nil, err
		}
		remain--

		topicFilterPairs = append(topicFilterPairs, &TopicFilterPair{topicFilter, qos})
	}
	return &SubscribePayload{topicFilterPairs}, nil
}

func extractLength(r *MQTTReader) (uint16, error) {
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

func extractQoS(r *MQTTReader) (uint8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	qos1 := refbit(b, 1)
	qos2 := refbit(b, 0)
	qos := qos2 << qos1
	return qos, nil
}
