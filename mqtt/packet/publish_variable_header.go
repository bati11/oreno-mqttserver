package packet

import (
	"bufio"
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

	_, err := r.ReadByte() // lengthMSB
	if err != nil {
		return PublishVariableHeader{}, err
	}

	lengthLSB, err := r.ReadByte()
	if err != nil {
		return PublishVariableHeader{}, err
	}
	if lengthLSB == 0 {
		return PublishVariableHeader{}, fmt.Errorf("length LSB should be > 0")
	}

	topicNameBytes := make([]byte, lengthLSB)
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
