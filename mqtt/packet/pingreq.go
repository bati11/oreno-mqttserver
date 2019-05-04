package packet

type Pingreq struct {
	FixedHeader *FixedHeader
}

func (reader *MQTTReader) ReadPingreq() (*Pingreq, error) {
	fixedHeader, err := reader.readFixedHeader()
	if err != nil {
		return nil, err
	}
	return &Pingreq{fixedHeader}, nil
}
