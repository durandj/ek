package web_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"testing"

	"github.com/durandj/ek/internal"
	"github.com/durandj/ek/internal/configuration"
	"github.com/durandj/ek/internal/logging"
	"github.com/durandj/ek/internal/sources"
	"github.com/durandj/ek/internal/web"
	"github.com/stretchr/testify/require"
)

func TestRedirectsListEndpoint(t *testing.T) {
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

	expectedRedirects := map[string]sources.Redirect{
		"test0": sources.Redirect{
			URLPattern: "http://example.com/test0",
		},
		"test1": sources.Redirect{
			URLPattern: "http://test1.example.com/",
		},
	}

	httpServer, err := internal.NewServer(internal.ServerParams{
		Config: config,
		Logger: logging.NewStructuredLoggerFromConfig(config),
		Source: sources.NewInMemorySource(sources.InMemorySourceParams{
			Redirects: expectedRedirects,
		}),
	})
	require.NoError(t, err)

	go func() {
		err := httpServer.Run(ctx)

		if ctx.Err() == nil {
			require.NoError(t, err, "HTTP server should start successfully")
		}
	}()

	httpClient := http.DefaultClient

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://"+httpServer.Address()+"/.api/redirects",
		nil,
	)
	require.NoError(t, err)

	response, err := httpClient.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	defer func() {
		_ = response.Body.Close()
	}()

	var actualRedirects web.RedirectsListResponse
	err = json.NewDecoder(response.Body).Decode(&actualRedirects)
	require.NoError(t, err)

	require.Equal(t, expectedRedirects, actualRedirects.Redirects)
}
