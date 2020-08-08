package clioptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	FlagDryRun         = "dry-run"
	FlagContext        = "context"
	FlagKubeConfigFile = "kubeconfig"
	FlagVerbosity      = "verbose"
)

func BindGlobalFlags(v *viper.Viper, flags *pflag.FlagSet) {

}
