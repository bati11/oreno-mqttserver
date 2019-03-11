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
	PublishFixedHeader
	ConnackVariableHeader
}

func (c Connack) ToBytes() []byte {
	var result []byte
	result = append(result, c.PublishFixedHeader.ToBytes()...)
	result = append(result, c.ConnackVariableHeader.ToBytes()...)
	return result
}

func newConnack() Connack {
	fixedHeader := PublishFixedHeader{
		PacketType:      2,
		RemainingLength: 2,
	}
	variableHeader := ConnackVariableHeader{SessionPresent: false}
	return Connack{fixedHeader, variableHeader}
}

func NewConnackForAccepted() Connack {
	result := newConnack()
	result.ReturnCode = 0
	return result
}

type ConnectError interface {
	Connack() Connack
	Error() string
}

type connectError struct {
	connack Connack
	msg     string
}

func (e connectError) Connack() Connack {
	return e.connack
}

func (e connectError) Error() string {
	return e.msg
}

func RefusedByUnacceptableProtocolVersion(s string) ConnectError {
	connack := newConnack()
	connack.ReturnCode = 1
	return connectError{connack, s}
}

func RefusedByIdentifierRejected(s string) ConnectError {
	connack := newConnack()
	connack.ReturnCode = 2
	return connectError{connack, s}
}
