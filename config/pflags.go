package config

import (
	"flag"

	"github.com/spf13/pflag"
)

func buildPFlagSet(fs *flag.FlagSet) (*pflag.FlagSet, error) {
	pfs := pflag.NewFlagSet(fs.Name(), pflag.ContinueOnError)
	pfs.AddGoFlagSet(fs)

	return pfs, nil
}
