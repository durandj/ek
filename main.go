package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/durandj/ek/internal"
	"github.com/durandj/ek/internal/configuration"
	"github.com/durandj/ek/internal/logging"
	"github.com/durandj/ek/internal/sources"
)

func Run() bool {
	config, err := configuration.NewConfigurationFromEnv()
	if err != nil {
		fmt.Printf("Unable to load configuration from environment: %v\n", err)

		return false
	}

	logger := logging.NewStructuredLoggerFromConfig(config)

	ctx := context.Background()
	ctx, cancelContext := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	// TODO: graceful shutdown timeout (20% of shutdown time [default of 30s])
	defer cancelContext()

	source, err := sources.NewFileSource(config.Source.FilePath)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"Unable to open configuration source",
			logging.Err(err),
		)

		return false
	}

	server, err := internal.NewServer(internal.ServerParams{
		Config: config,
		Logger: logger,
		Source: source,
	})
	if err != nil {
		logger.ErrorContext(ctx, "Unable to create server", logging.Err(err))
	}

	if err := server.Run(ctx); err != nil {
		logger.ErrorContext(ctx, "Fatal error in server", logging.Err(err))

		return false
	}

	return true
}

func main() {
	if ok := Run(); !ok {
		os.Exit(1)
	}
}
