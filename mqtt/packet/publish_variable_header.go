package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type PublishVariableHeader struct {
	TopicName        string
	PacketIdentifier *uint16
}

func ToPublishVariableHeader(fixedHeader FixedHeader, bs []byte) (PublishVariableHeader, []byte, error) {
	if fixedHeader.PacketType != PUBLISH {
		return PublishVariableHeader{}, nil, fmt.Errorf("packet type is invalid. it got is %v", fixedHeader.PacketType)
	}

	if len(bs) < 2 {
		return PublishVariableHeader{}, nil, fmt.Errorf("len(bs) should be >= 2")
	}

	lengthMSB := bs[0]
	lengthLSB := bs[1]
	if lengthLSB == 0 {
		return PublishVariableHeader{}, nil, fmt.Errorf("length LSB should be > 0")
	}
	if len(bs) < (2 + int(lengthLSB)) {
		return PublishVariableHeader{}, nil, fmt.Errorf("len(bs) should be >= 2+LSB")
	}
	topicName := bs[2+lengthMSB : 2+lengthLSB]
	if strings.ContainsAny(string(topicName), "# & +") {
		return PublishVariableHeader{}, nil, fmt.Errorf("topic name must not contain wildcard. it got is %v", topicName)
	}

	var packetIdentifier uint16
	if fixedHeader.QoS == QoS1 || fixedHeader.QoS == QoS2 {
		packetIdentifierMSB := bs[2+lengthLSB]
		packetIdentifierLSB := bs[2+lengthLSB+1]
		binary.Read(bytes.NewBuffer([]byte{packetIdentifierMSB, packetIdentifierLSB}), binary.BigEndian, &packetIdentifier)
		result := PublishVariableHeader{string(topicName), &packetIdentifier}
		return result, bs[2+lengthLSB+2:], nil
	}
	result := PublishVariableHeader{string(topicName), nil}
	return result, bs[2+lengthLSB:], nil
}
