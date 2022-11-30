package peer

type Network interface {
	Pong(nodeId int) (string, error)
}
