package sender

type ExternalSender interface {
	Send(msg string, nodeIds []int) []int
	Gossip(msg string, numPeersToSend int) []int
}
