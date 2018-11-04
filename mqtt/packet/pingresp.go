package packet

type Pingresp struct {
	FixedHeader
}

func NewPingresp() Pingresp {
	fixedHeader := FixedHeader{
		PacketType:      PINGRESP,
		RemainingLength: 0,
	}
	return Pingresp{fixedHeader}
}

func (c *Pingresp) ToBytes() []byte {
	return c.FixedHeader.ToBytes()
}
