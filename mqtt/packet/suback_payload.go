package packet

type SubackPayload struct {
	ReturnCode byte
}

func NewSubackPayload(qos QoS) SubackPayload {
	switch qos {
	case QoS0:
		return SubackPayload{0x00}
	case QoS1:
		return SubackPayload{0x01}
	case QoS2:
		return SubackPayload{0x02}
	}
	return NewSubackPayloadForFail()
}

func NewSubackPayloadForFail() SubackPayload {
	return SubackPayload{0x80}
}

func (p *SubackPayload) ToBytes() []byte {
	return []byte{p.ReturnCode}
}
