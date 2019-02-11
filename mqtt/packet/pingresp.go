package packet

type Pingresp struct {
	FixedHeader
}

func (p *Pingresp) ToBytes() []byte {
	var result []byte
	result = append(result, p.FixedHeader.ToBytes()...)
	return result
}

func NewPingresp() Pingresp {
	fixedHeader := FixedHeader{
		PacketType:      PINGRESP,
		RemainingLength: 0,
	}
	return Pingresp{fixedHeader}
}
