package cmd

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func RegisterRootFlags(flags *pflag.FlagSet) {
	var err error
	if err = viper.BindPFlag("m", flags.Lookup("m")); err != nil {
		panic(err)
	}
}
