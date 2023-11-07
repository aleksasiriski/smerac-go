package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type Flags struct {
	Globals

	// flags
	ConfigDirPath string `type:"path" default:"${config_folder}" env:"SMERAC_CONFIG_DIR" help:"Data folder path"`
	LogDirPath    string `type:"path" default:"${log_folder}" env:"SMERAC_LOG_DIR" help:"Log folder path"`
	Verbosity     int8   `type:"counter" default:"0" short:"v" env:"SMERAC_VERBOSITY" help:"Log level verbosity"`
}

var (
	// release variables
	Version   string
	Timestamp string
	GitCommit string
)

type Globals struct {
	Version versionFlag `name:"version" help:"Print version information and quit"`
}

type versionFlag string

func (v versionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v versionFlag) IsBool() bool                         { return true }
func (v versionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}
