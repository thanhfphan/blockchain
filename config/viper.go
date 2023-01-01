package config

import (
	"flag"

	"github.com/spf13/viper"
)

func BuildViper(fs *flag.FlagSet, args []string) (*viper.Viper, error) {
	pfs, err := buildPFlagSet(fs)
	if err != nil {
		return nil, err
	}
	if err := pfs.Parse(args); err != nil {
		return nil, err
	}

	v := viper.New()

	if err := v.BindPFlags(pfs); err != nil {
		return nil, err
	}
	return v, nil
}
