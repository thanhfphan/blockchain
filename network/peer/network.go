package peer

type Network interface {
	Pong(nodeID int) (string, error)

	Connected(nodeID int)
	Disconnected(nodeID int)
}
