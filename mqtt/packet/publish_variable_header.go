package packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type PublishVariableHeader struct {
	LengthMSB        byte
	LengthLSB        byte
	TopicName        string
	PacketIdentifier *uint16
}

func NewPublishVariableHeader(topicName string) *PublishVariableHeader {
	b := make([]byte, 2)
	length := len(topicName)
	// TODO check length value

	binary.BigEndian.PutUint16(b, uint16(length))
	return &PublishVariableHeader{
		LengthMSB: b[0],
		LengthLSB: b[1],
		TopicName: topicName,
	}
}

func (p *PublishVariableHeader) Length() uint {
	result := uint(len(p.TopicName) + 2)
	if p.PacketIdentifier != nil {
		return result + uint(*p.PacketIdentifier)
	}
	return result
}

func (p *PublishVariableHeader) ToBytes() []byte {
	var result []byte
	result = append(result, p.LengthMSB)
	result = append(result, p.LengthLSB)
	result = append(result, []byte(p.TopicName)...)

	if p.PacketIdentifier != nil {
		bs := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(bs, *p.PacketIdentifier)
		result = append(result, bs...)
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

	result := PublishVariableHeader{lengthMSB, lengthLSB, string(topicNameBytes), nil}
	return &result, nil
}
