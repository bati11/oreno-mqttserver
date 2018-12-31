package packet

type ConnackVariableHeader struct {
	SessionPresent bool
	ReturnCode     uint8
}

func (h ConnackVariableHeader) ToBytes() []byte {
	var result []byte
	if h.SessionPresent {
		result = append(result, 1)
	} else {
		result = append(result, 0)
	}
	result = append(result, h.ReturnCode)
	return result
}

type Connack struct {
	FixedHeader
	ConnackVariableHeader
}

func (c Connack) ToBytes() []byte {
	var result []byte
	result = append(result, c.FixedHeader.ToBytes()...)
	result = append(result, c.ConnackVariableHeader.ToBytes()...)
	return result
}

func NewConnackForAccepted() Connack {
	result := newConnack()
	result.ReturnCode = 0
	return result
}

func NewConnackForRefusedByUnacceptableProtocolVersion() Connack {
	result := newConnack()
	result.ReturnCode = 1
	return result
}

func NewConnackForRefusedByIdentifierRejected() Connack {
	result := newConnack()
	result.ReturnCode = 2
	return result
}

func newConnack() Connack {
	fixedHeader := FixedHeader{
		PacketType:      2,
		RemainingLength: 2,
	}
	variableHeader := ConnackVariableHeader{SessionPresent: false}
	return Connack{fixedHeader, variableHeader}
}
