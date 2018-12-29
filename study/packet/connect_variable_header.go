package packet

type ConnectFlags struct {
	CleanSession bool
	WillFlag     bool
	WillQoS      uint8
	WillRetain   bool
	PasswordFlag bool
	UserNameFlag bool
}

type ConnectVariableHeader struct {
	ProtocolName  string
	ProtocolLevel uint8
	ConnectFlags  ConnectFlags
	KeepAlive     uint16
}

func ToConnectVariableHeader(fixedHeader FixedHeader, bs []byte) (ConnectVariableHeader, error) {
	return ConnectVariableHeader{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		ConnectFlags:  ConnectFlags{UserNameFlag: true, PasswordFlag: true, WillRetain: false, WillQoS: 1, WillFlag: true, CleanSession: true},
		KeepAlive:     10,
	}, nil
}
