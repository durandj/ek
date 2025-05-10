package web_test

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/durandj/ek/internal"
	"github.com/durandj/ek/internal/configuration"
	"github.com/durandj/ek/internal/logging"
	"github.com/durandj/ek/internal/sources"
	"github.com/stretchr/testify/require"
)

func TestHealthcheckEndpoint(t *testing.T) {
	ctx := t.Context()

	config := configuration.Configuration{
		Environment: configuration.EnvironmentDevelopment,
		Logging: configuration.LoggingConfiguration{
			Level: configuration.LoggingLevel(slog.LevelInfo),
		},
		HTTP: configuration.HTTPConfiguration{
			Interface: "127.0.0.1",
			Port:      0,
		},
		Source: configuration.SourceConfiguration{},
	}

	httpServer, err := internal.NewServer(internal.ServerParams{
		Config: config,
		Logger: logging.NewStructuredLoggerFromConfig(config),
		Source: sources.NewInMemorySource(sources.InMemorySourceParams{
			Redirects: nil,
		}),
	})
	require.NoError(t, err)

	go func() {
		err := httpServer.Run(ctx)
		require.NoError(t, err, "HTTP server should start successfully")
	}()

	httpClient := http.DefaultClient

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://"+httpServer.Address()+"/.api/healthcheck",
		nil,
	)
	require.NoError(t, err)

	response, err := httpClient.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)
}
