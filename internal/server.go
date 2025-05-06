package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/durandj/ek/internal/configuration"
	"github.com/durandj/ek/internal/logging"
	"github.com/durandj/ek/internal/sources"
	"github.com/durandj/ek/internal/web"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

type ServerParams struct {
	Config configuration.Configuration
	Logger *slog.Logger
	Source sources.Source
}

func NewServer(params ServerParams) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: params.Config.HTTP.Address(),
			Handler: web.NewRouter(web.RouterParams{
				Config: params.Config,
				Logger: params.Logger,
				Source: params.Source,
			}),
			ReadTimeout:                  0,
			ReadHeaderTimeout:            readHeaderTimeout,
			WriteTimeout:                 0,
			IdleTimeout:                  time.Minute,
			MaxHeaderBytes:               0,
			ConnState:                    nil,
			BaseContext:                  nil,
			ConnContext:                  nil,
			Protocols:                    nil,
			TLSConfig:                    nil,
			TLSNextProto:                 nil,
			HTTP2:                        nil,
			DisableGeneralOptionsHandler: false,
			ErrorLog:                     nil, // TODO: structured logger?
		},
		logger: params.Logger,
	}
}

func (server *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		// TODO: this should probably be a timeout context instead so we can force terminate
		if err := server.httpServer.Shutdown(context.Background()); err != nil {
			server.logger.ErrorContext(
				ctx,
				"Error when shutting down HTTP server",
				logging.Err(err),
			)
		}
	}()
	server.logger.Info(
		"Starting HTTP server",
		slog.String("address", server.httpServer.Addr),
	)
	if err := server.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("error in HTTP server: %w", err)
	}

	return nil
}

const readHeaderTimeout = 10 * time.Second
