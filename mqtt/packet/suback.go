package packet

type Suback struct {
	FixedHeader
	SubackVariableHeader
	SubackPayload
}

func NewSuback(subscribeVariableHeader SubscribeVariableHeader, subscribePayload SubscribePayload) Suback {
	variableHeader := SubackVariableHeader{subscribeVariableHeader.PacketIdentifier}
	payload := NewSubackPayload(subscribePayload.QoS)

	fixedHeader := FixedHeader{
		PacketType:      SUBACK,
		RemainingLength: uint(len(variableHeader.ToBytes())) + uint(len(payload.ToBytes())),
	}
	return Suback{fixedHeader, variableHeader, payload}
}

func (c *Suback) ToBytes() []byte {
	var result []byte
	result = append(result, c.FixedHeader.ToBytes()...)
	result = append(result, c.SubackVariableHeader.ToBytes()...)
	result = append(result, c.SubackPayload.ToBytes()...)

	return result
}
