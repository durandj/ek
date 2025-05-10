package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/durandj/ek/internal/configuration"
	"github.com/durandj/ek/internal/logging"
	"github.com/durandj/ek/internal/sources"
	"github.com/durandj/ek/internal/web"
)

type Server struct {
	httpListener net.Listener
	httpServer   *http.Server
	logger       *slog.Logger
}

type ServerParams struct {
	Config configuration.Configuration
	Logger *slog.Logger
	Source sources.Source
}

func NewServer(params ServerParams) (*Server, error) {
	listener, err := net.Listen("tcp", params.Config.HTTP.Address())
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP listener: %w", err)
	}

	return &Server{
		httpListener: listener,
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
	}, nil
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
	if err := server.httpServer.Serve(server.httpListener); err != nil {
		return fmt.Errorf("error in HTTP server: %w", err)
	}

	return nil
}

func (server *Server) Address() string {
	return server.httpListener.Addr().String()
}

const readHeaderTimeout = 10 * time.Second
