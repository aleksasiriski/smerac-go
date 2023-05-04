package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/bwmarrin/discordgo"
)

var (
	// CLI
	cli struct {
		// flags
		Config    string `type:"path" default:"${config_dir}" env:"SMERAC_GO_CONFIG" help:"Config file path"`
		Log       string `type:"path" default:"${log_file}" env:"SMERAC_GO_LOG" help:"Log file path"`
		Verbosity int    `type:"counter" default:"0" short:"v" env:"SMERAC_GO_VERBOSITY" help:"Log level verbosity"`
	}
)

func main() {
	// parse cli
	ctx := kong.Parse(&cli,
		kong.Name("smerac-go"),
		kong.Description("Discord bot that allows users to choose their own role and pull calendar from Google."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Summary: true,
			Compact: true,
		}),
		kong.Vars{
			"config_dir": "/etc/smerac-go",
			"log_file":   "/var/log/smerac-go.log",
		},
	)

	if err := ctx.Validate(); err != nil {
		fmt.Println("Failed parsing cli:", err)
		os.Exit(1)
	}

	// logger
	logger := log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out:        os.Stderr,
	}, zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out: &lumberjack.Logger{
			Filename:   cli.Log,
			MaxSize:    5,
			MaxAge:     14,
			MaxBackups: 5,
		},
		NoColor: true,
	}))

	switch {
	case cli.Verbosity == 1:
		log.Logger = logger.Level(zerolog.DebugLevel)
	case cli.Verbosity > 1:
		log.Logger = logger.Level(zerolog.TraceLevel)
	default:
		log.Logger = logger.Level(zerolog.InfoLevel)
	}

	// config
	config, err := LoadConfig(cli.Config)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Cannot load config:")
	}

	// discordgo
	discord, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed initialising discord API:")
	}

	// discord connection
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info().
			Msg(fmt.Sprintf("Logged in as: %s#%s", s.State.User.Username, s.State.User.Discriminator))
	})
	err = discord.Open()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed connecting to discord servers:")
	}

	// smerac-go
	var worker conc.WaitGroup
	worker.Go(func() {
		worker.Go(func() {
			err := setupSlashCommands(config.Roles, discord)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Failed setting up slash commands:")
			}
		})
		worker.Go(func() {
			err := updateCalendars(config.Calendars, discord, config.Google)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Failed updating calendars:")
			}
		})
	})

	// handle interrupt signal
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, os.Kill)
	<-quitChannel
}
