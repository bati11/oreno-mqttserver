package packet

type Pingresp struct {
	DefaultFixedHeader
}

func (p *Pingresp) ToBytes() []byte {
	var result []byte
	result = append(result, p.DefaultFixedHeader.ToBytes()...)
	return result
}

func NewPingresp() Pingresp {
	fixedHeader := DefaultFixedHeader{
		PacketType:      PINGRESP,
		RemainingLength: 0,
	}
	return Pingresp{fixedHeader}
}
