package packet

import (
	"bufio"
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

func ToPublishVariableHeader(fixedHeader FixedHeader, r *bufio.Reader) (PublishVariableHeader, error) {
	if fixedHeader.PacketType != 3 {
		return PublishVariableHeader{}, fmt.Errorf("packet type is invalid. it got is %v", fixedHeader.PacketType)
	}

	lengthMSB, err := r.ReadByte()
	if err != nil {
		return PublishVariableHeader{}, err
	}
	lengthLSB, err := r.ReadByte()
	if err != nil {
		return PublishVariableHeader{}, err
	}
	n := binary.BigEndian.Uint16([]byte{lengthMSB, lengthLSB})
	if n == 0 {
		return PublishVariableHeader{}, fmt.Errorf("topic name length should be > 0")
	}

	topicNameBytes := make([]byte, n)
	_, err = io.ReadFull(r, topicNameBytes)
	if err != nil {
		return PublishVariableHeader{}, err
	}

	topicName := string(topicNameBytes)
	if strings.ContainsAny(topicName, "# +") {
		return PublishVariableHeader{}, fmt.Errorf("topic name must not contain wildcard. it got is %v", topicName)
	}

	result := PublishVariableHeader{string(topicNameBytes), nil}
	return result, nil
}
