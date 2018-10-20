package packet

type ConnackVariableHeader struct {
	SessionPresent bool
	ReturnCode     uint8
}

func (h *ConnackVariableHeader) ToBytes() []byte {
	var result []byte
	if h.SessionPresent {
		result = append(result, 1)
	} else {
		result = append(result, 0)
	}
	result = append(result, h.ReturnCode)
	return result
}
