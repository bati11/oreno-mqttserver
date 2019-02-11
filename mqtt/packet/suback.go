package packet

type Suback struct {
	FixedHeader
	SubackVariableHeader
	SubackPayload
}

func NewSubackForSuccess(packetIdentifier uint16, qoss []uint8) Suback {
	variableHeader := SubackVariableHeader{packetIdentifier}
	payload := SubackPayload{qoss}
	fixedHeader := FixedHeader{
		PacketType:      SUBACK,
		RemainingLength: variableHeader.Length() + payload.Length(),
	}
	return Suback{fixedHeader, variableHeader, payload}
}

func (s *Suback) ToBytes() []byte {
	var result []byte
	result = append(result, s.FixedHeader.ToBytes()...)
	result = append(result, s.SubackVariableHeader.ToBytes()...)
	result = append(result, s.SubackPayload.ToBytes()...)
	return result
}
