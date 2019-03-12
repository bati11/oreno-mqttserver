package packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type PublishVariableHeader struct {
	TopicName        string
	PacketIdentifier *uint16
}

func (p *PublishVariableHeader) Length() uint {
	result := uint(len(p.TopicName) + 2)
	if p.PacketIdentifier != nil {
		return result + uint(*p.PacketIdentifier)
	}
	return result
}

func (reader *MQTTReader) readPublishVariableHeader() (*PublishVariableHeader, error) {
	lengthMSB, err := reader.r.ReadByte()
	if err != nil {
		return nil, err
	}
	lengthLSB, err := reader.r.ReadByte()
	if err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint16([]byte{lengthMSB, lengthLSB})
	if n == 0 {
		return nil, fmt.Errorf("topic name length should be > 0")
	}

	topicNameBytes := make([]byte, n)
	_, err = io.ReadFull(reader.r, topicNameBytes)
	if err != nil {
		return nil, err
	}

	topicName := string(topicNameBytes)
	if strings.ContainsAny(topicName, "# +") {
		return nil, fmt.Errorf("topic name must not contain wildcard. it got is %v", topicName)
	}

	result := PublishVariableHeader{string(topicNameBytes), nil}
	return &result, nil
}
