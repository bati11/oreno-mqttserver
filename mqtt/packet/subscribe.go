package packet

type Subscribe struct {
	FixedHeader    *FixedHeader
	VariableHeader *SubscribeVariableHeader
	Payload        *SubscribePayload
}

func (reader *MQTTReader) ReadSubscribe() (*Subscribe, error) {
	fixedHeader, err := reader.readFixedHeader()
	if err != nil {
		return nil, err
	}
	variableHeader, err := reader.readSubscribeVariableHeader()
	if err != nil {
		return nil, err
	}
	payloadLength := fixedHeader.RemainingLength - variableHeader.Length()
	payload, err := reader.readSubscribePayload(payloadLength)
	if err != nil {
		return nil, err
	}
	return &Subscribe{fixedHeader, variableHeader, payload}, nil
}
