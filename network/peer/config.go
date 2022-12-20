package peer

import (
	"time"

	"github.com/thanhfphan/blockchain/message"
)

type Config struct {
	Network        Network
	MessageCreator message.Creator
	PongTimeout    time.Duration
}
