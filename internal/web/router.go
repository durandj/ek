package web

import (
	"log/slog"
	"time"

	"github.com/durandj/ek/internal/configuration"
	"github.com/durandj/ek/internal/sources"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

type RouterParams struct {
	Config configuration.Configuration
	Logger *slog.Logger
	Source sources.Source
}

func NewRouter(params RouterParams) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		// The logger needs to be the first middleware
		httplog.RequestLogger(newRequestLogger(params.Config)),
		middleware.Recoverer,
		middleware.RealIP,
		middleware.RedirectSlashes,
		middleware.CleanPath,
		middleware.Timeout(1*time.Minute),
		middleware.ContentCharset("UTF-8", "Latin-1", ""),
	)

	// TODO: rate limits

	router.Route("/.api", func(apiRouter chi.Router) {
		apiRouter.Get("/healthcheck", HealthcheckEndpoint())

		apiRouter.Route("/redirects", func(redirectsRouter chi.Router) {
			redirectsRouter.Get(
				"/",
				RedirectsListEndpoint(RedirectsListEndpointParams{
					Logger: params.Logger,
					Source: params.Source,
				}),
			)
		})
	})

	redirectParams := RedirectEndpointParams{
		Logger: params.Logger,
		Source: params.Source,
	}
	router.Get("/{key}", RedirectEndpoint(redirectParams))
	router.Get("/{key}/*", RedirectEndpoint(redirectParams))

	return router
}

func newRequestLogger(config configuration.Configuration) *httplog.Logger {
	isDevEnv := config.Environment == configuration.EnvironmentDevelopment

	return httplog.NewLogger("datumo-api", httplog.Options{
		JSON:               !isDevEnv,
		LogLevel:           config.Logging.Level.ToSlogLevel(),
		LevelFieldName:     "level",
		Concise:            true,
		HideRequestHeaders: []string{},
		RequestHeaders:     isDevEnv,
		ResponseHeaders:    isDevEnv,
		MessageFieldName:   "message",
		TimeFieldName:      "timestamp",
		TimeFieldFormat:    time.RFC3339,
		SourceFieldName:    "source",
		Tags: map[string]string{
			// TODO: version
			"env": config.Environment.String(),
		},
		ReplaceAttrsOverride: nil,
		QuietDownRoutes: []string{
			"/healthcheck",
		},
		QuietDownPeriod: quietDownPeriod,
		Writer:          nil,
		Trace:           nil,
	})
}

const quietDownPeriod = 10 * time.Second
