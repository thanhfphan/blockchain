package config

import (
	"flag"
	"time"

	"github.com/thanhfphan/blockchain/utils/constants"
)

const (
	DefaultStakingPort     = 4012
	DefaultStakingIPAdress = "127.0.1"
)

func BuildFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(constants.AppName, flag.ContinueOnError)
	addNodeFlags(fs)

	return fs
}

func addNodeFlags(fs *flag.FlagSet) {
	fs.Duration(NetworkPingFrequencyKey, 30*time.Second, "Frequency of pinging other peers")
	fs.Duration(NetworkPingTimeoutKey, 10*time.Second, "Timeout value for Ping-Pong with a peer")

	fs.String(PublicIPKey, DefaultStakingIPAdress, "IP of this node for P2P communication")
	fs.Uint(StakingPortKey, DefaultStakingPort, "Port of the consensus server")
}
