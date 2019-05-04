package packet

import (
	"io"
)

type Publish struct {
	FixedHeader    *PublishFixedHeader
	VariableHeader *PublishVariableHeader
	Payload        PublishPayload
}

func (p Publish) ToBytes() []byte {
	var result []byte
	result = append(result, p.FixedHeader.ToBytes()...)
	result = append(result, p.VariableHeader.ToBytes()...)
	result = append(result, []byte(string(p.Payload))...)
	return result
}

type PublishPayload []byte

func (p *Publish) PayloadLength() uint {
	return p.FixedHeader.RemainingLength - p.VariableHeader.Length()
}

func NewPublish(topicName string, message []byte) *Publish {
	variableHeader := NewPublishVariableHeader(topicName)
	fixedHeader := NewPublishFixedHeader(PUBLISH, variableHeader.Length()+uint(len(message)))
	return &Publish{fixedHeader, variableHeader, message}
}

func (reader *MQTTReader) ReadPublish() (*Publish, error) {
	fixedHeader, err := reader.readPublishFixedHeader()
	if err != nil {
		return nil, err
	}
	variableHeader, err := reader.readPublishVariableHeader()
	if err != nil {
		return nil, err
	}
	payloadLength := fixedHeader.RemainingLength - variableHeader.Length()
	payload := make([]byte, payloadLength)
	_, err = io.ReadFull(reader.r, payload)
	if err != nil {
		return nil, err
	}
	return &Publish{fixedHeader, variableHeader, payload}, nil
}
