package clioptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// FlagDryRun is for --dry-run
	FlagDryRun = "dry-run"
	// FlagContext is for --context passed to kubectl
	FlagContext = "context"
	// FlagKubeConfigFile is for --kubeconfig passed to kubectl
	FlagKubeConfigFile = "kubeconfig"
	// FlagVerbosity is for --verbose
	FlagVerbosity = "verbose"
)

// BindGlobalFlags binds global (reusable flags) to a particular flagset, such as from a subcommand
func BindGlobalFlags(v *viper.Viper, flags *pflag.FlagSet) {
}
