package message

type MessageType string

const (
	MessageTypePing     MessageType = "ping"
	MessageTypePong     MessageType = "pong"
	MessageTypeHello    MessageType = "hello"
	MessageTypePeerList MessageType = "peer-list"
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
	HelloTime  uint64 `json:"hello_time"`
	IPAddress  []byte `json:"ip_address"`
	IPPort     uint16 `json:"ip_port"`
	SignedTime uint64 `json:"signed_time:`
	Signature  []byte `json:"signature"`
}

type MessagePeerList struct {
	PeerList []*PeerList `json:"peer_list"`
}

type PeerList struct {
	Cert      []byte `json:"cert"`
	IPAddress []byte `json:"ip_address"`
	IPPort    uint16 `json:"ip_port"`
	Signature []byte `json:"signature"`
	TxID      []byte `json:"tx_id"`
	Timestamp uint64 `json:"timestamp"`
}
