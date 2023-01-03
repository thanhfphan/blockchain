package genesis

type GenesisNode struct {
	IP     string
	NodeID string
}

func GetGenesisNodes() []*GenesisNode {
	gNodes := []*GenesisNode{
		{
			IP:     "127.0.0.1:6001",
			NodeID: "NodeID-7ctjx7KwntsshR3AhyPK6YtAUqoV8JEuQWZJ48Nyv4sk",
		},
	}

	return gNodes
}
