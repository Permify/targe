package users

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func RegisterUsersFlags(flags *pflag.FlagSet) {
	var err error

	// Config File
	if err = viper.BindPFlag("config.file", flags.Lookup("config")); err != nil {
		panic(err)
	}

	if err = viper.BindPFlag("user", flags.Lookup("user")); err != nil {
		panic(err)
	}
	if err = viper.BindEnv("user", "KIVO_USER"); err != nil {
		panic(err)
	}

	if err = viper.BindPFlag("action", flags.Lookup("action")); err != nil {
		panic(err)
	}
	if err = viper.BindEnv("action", "KIVO_ACTION"); err != nil {
		panic(err)
	}

	if err = viper.BindPFlag("policy", flags.Lookup("policy")); err != nil {
		panic(err)
	}
	if err = viper.BindEnv("policy", "KIVO_POLICY"); err != nil {
		panic(err)
	}

	if err = viper.BindPFlag("resource", flags.Lookup("resource")); err != nil {
		panic(err)
	}
	if err = viper.BindEnv("resource", "KIVO_RESOURCE"); err != nil {
		panic(err)
	}

	if err = viper.BindPFlag("service", flags.Lookup("service")); err != nil {
		panic(err)
	}
	if err = viper.BindEnv("service", "KIVO_SERVICE"); err != nil {
		panic(err)
	}

	if err = viper.BindPFlag("policy_option", flags.Lookup("policy-option")); err != nil {
		panic(err)
	}
	if err = viper.BindEnv("policy_option", "KIVO_POLICY_OPTION"); err != nil {
		panic(err)
	}
}
