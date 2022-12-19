package message

type MessageType string

const (
	MessageTypePing  MessageType = "ping"
	MessageTypePong  MessageType = "pong"
	MessageTypeHello MessageType = "hello"
)

type Message struct {
	Type    MessageType `json:"type"`
	Message any
}

type MessagePing struct {
	Message string `json:"message"`
}

type MessagePong struct {
	Message string `json:"message"`
}

type MessageHello struct {
	Message string `json:"message"`
	MyTime  uint64 `json:"my_time"`
}
