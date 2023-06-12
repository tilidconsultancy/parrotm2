package domain

type (
	IncomingMessage struct {
		Id      string
		Message interface{}
	}
)

func NewIncomingMessage(id string, msg interface{}) *IncomingMessage {
	return &IncomingMessage{
		Id:      id,
		Message: msg,
	}
}
