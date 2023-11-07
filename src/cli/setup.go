package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

func Setup() Flags {
	var cli Flags
	ctx := kong.Parse(&cli,
		kong.Name("brzaguza"),
		kong.Description("Fastasst metasearch engine"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Summary: true,
			Compact: true,
		}),
		kong.Vars{
			"version":       fmt.Sprintf("%v (%v@%v)", Version, GitCommit, Timestamp),
			"config_folder": ".",
			"log_folder":    "./log",
		},
	)

	if err := ctx.Validate(); err != nil {
		log.Panic().Err(err).Msg("Failed parsing cli")
	}

	return cli
}
