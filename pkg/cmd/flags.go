package cmd

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func RegisterRootFlags(flags *pflag.FlagSet) {
	var err error
	if err = viper.BindPFlag("ticket", flags.Lookup("ticket")); err != nil {
		panic(err)
	}
}
