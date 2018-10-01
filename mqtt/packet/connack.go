package packet

type Connack struct {
	FixedHeader
	ConnackVariableHeader
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
		PacketType:      CONNACK,
		RemainingLength: 2,
	}
	variableHeader := ConnackVariableHeader{SessionPresent: false}
	return Connack{fixedHeader, variableHeader}
}

func (c *Connack) ToBytes() []byte {
	var result []byte
	result = append(result, c.FixedHeader.ToBytes()...)
	result = append(result, c.ConnackVariableHeader.ToBytes()...)
	return result
}
