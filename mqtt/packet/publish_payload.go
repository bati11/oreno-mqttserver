package packet

type PublishPayload struct {
	Message []byte
}

func ToPublishPayload(bs []byte) PublishPayload {
	return PublishPayload{bs}
}
