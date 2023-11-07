package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aleksasiriski/smerac-go/src/calendar"
	"github.com/aleksasiriski/smerac-go/src/cli"
	"github.com/aleksasiriski/smerac-go/src/config"
	"github.com/aleksasiriski/smerac-go/src/logger"
)

func main() {
	// parse cli arguments
	cliFlags := cli.Setup()

	// signal interrupt (CTRL+C)
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// configure logging
	logger.Setup(cliFlags.LogDirPath, cliFlags.Verbosity)

	// load config file
	conf := config.New()
	conf.Load(cliFlags.ConfigDirPath, cliFlags.LogDirPath)

	// startup
	calendar.Update(ctx, conf)
}
