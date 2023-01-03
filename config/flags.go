package config

import (
	"flag"
	"time"

	"github.com/thanhfphan/blockchain/utils/constants"
)

const (
	DefaultStakingPort     = 4012
	DefaultStakingIPAdress = "127.0.0.1"

	DefaultPingFrequency = 30 * time.Second
	DefaultPingTimeout   = 10 * time.Second
	DefaultDialerTimeout = 10 * time.Second
)

func BuildFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(constants.AppName, flag.ContinueOnError)
	addNodeFlags(fs)

	return fs
}

func addNodeFlags(fs *flag.FlagSet) {
	fs.Duration(NetworkPingFrequencyKey, DefaultPingFrequency, "Frequency of pinging other peers")
	fs.Duration(NetworkPingTimeoutKey, DefaultPingTimeout, "Timeout value for Ping-Pong with a peer")
	fs.Duration(NetworkDialerTimeoutKey, DefaultDialerTimeout, "Timeout value for dial with other peers")

	// fs.String(PublicIPKey, DefaultStakingIPAdress, "IP of this node for P2P communication") //only support local for now
	fs.Uint(StakingPortKey, DefaultStakingPort, "Port of the consensus server")
	// fs.String(StakingTLSKeyPathKey, "", "Path to the TLS private key. If not specified, create a random one")
	// fs.String(StakingTLSCertPathKey, "", "Path to the TLS certificate key. If not specified, create a random one")
}
