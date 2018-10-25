package packet

type MQTTPacket interface {
	ToBytes() []byte
}
