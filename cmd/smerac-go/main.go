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
		Config    string `type:"path" default:"${config_dir}" env:"RFFMPEG_AUTOSCALER_CONFIG" help:"Config file path"`
		Log       string `type:"path" default:"${log_file}" env:"RFFMPEG_AUTOSCALER_LOG" help:"Log file path"`
		Verbosity int    `type:"counter" default:"0" short:"v" env:"RFFMPEG_AUTOSCALER_VERBOSITY" help:"Log level verbosity"`
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
			"config_dir": "/config",
			"log_file":   "/config/log/smerac.log",
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

	// display initialised banner
	log.Info().
		Str("Discord", "success").
		Msg("Initialised")

	// smerac-go
	var helper conc.WaitGroup
	var worker conc.WaitGroup
	helper.Go(func() {
		for {
			worker.Go(func() {
				err := fmt.Errorf("woohoo")
				if err != nil {
					log.Error().
						Err(err).
						Msg("Failed while doing:")
				}
			})
			time.Sleep(time.Minute * 5)
		}
	})

	// handle interrupt signal
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	worker.Wait()

	// testing
	fmt.Println(config)
	fmt.Println(discord)
}
