package packet

type SubackPayload struct {
	ReturnCodes []byte
}

func (s *SubackPayload) Length() uint {
	return uint(len(s.ReturnCodes))
}

func (s *SubackPayload) ToBytes() []byte {
	return s.ReturnCodes
}
