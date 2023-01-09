package config

import (
	"flag"
	"time"

	"github.com/thanhfphan/blockchain/utils/constants"
	"github.com/thanhfphan/blockchain/utils/logging"
)

const (
	DefaultStakingPort           = 4012
	DefaultStakingIPAdress       = "127.0.0.1"
	DefaultNumberValidatorOfPeer = 10

	DefaultPingTimeout   = 10 * time.Second
	DefaultPingFrequency = 3 * DefaultPingTimeout / 4
	DefaultDialerTimeout = 10 * time.Second

	DefaultLogLevel = logging.VerboStr
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
	fs.Uint(NetworkNumberValidatorOfPeerKey, DefaultNumberValidatorOfPeer, "Number validator of a peer")

	fs.String(PublicIPKey, DefaultStakingIPAdress, "IP of this node for P2P communication") //only support local for now
	fs.Uint(StakingPortKey, DefaultStakingPort, "Port of the consensus server")
	// fs.String(StakingTLSKeyPathKey, "", "Path to the TLS private key. If not specified, create a random one")
	// fs.String(StakingTLSCertPathKey, "", "Path to the TLS certificate key. If not specified, create a random one")

	fs.String(LogLevelKey, DefaultLogLevel, "Log level")
}
