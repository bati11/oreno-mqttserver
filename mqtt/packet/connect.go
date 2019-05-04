package packet

type Connect struct {
	FixedHeader    *FixedHeader
	VariableHeader *ConnectVariableHeader
	Payload        *ConnectPayload
}

func (reader *MQTTReader) ReadConnect() (*Connect, error) {
	fixedHeader, err := reader.readFixedHeader()
	if err != nil {
		return nil, err
	}
	variableHeader, err := reader.readConnectVariableHeader()
	if err != nil {
		return nil, err
	}
	payload, err := reader.readConnectPayload()
	if err != nil {
		return nil, err
	}
	return &Connect{fixedHeader, variableHeader, payload}, nil
}
