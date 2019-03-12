package packet

import (
	"io"
)

type Publish struct {
	FixedHeader    *FixedHeader
	VariableHeader *PublishVariableHeader
	Payload        PublishPayload
}

type PublishPayload []byte

func (p *Publish) PayloadLength() uint {
	return p.FixedHeader.RemainingLength - p.VariableHeader.Length()
}

func (reader *MQTTReader) ReadPublish() (*Publish, error) {
	fixedHeader, err := reader.readFixedHeader()
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
