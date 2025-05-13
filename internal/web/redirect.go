package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/durandj/ek/internal/logging"
	"github.com/durandj/ek/internal/sources"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type RedirectEndpointParams struct {
	Logger *slog.Logger
	Source sources.Source
}

func RedirectEndpoint(params RedirectEndpointParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		redirectKey := chi.URLParam(r, "key")
		redirectContextPath := chi.URLParam(r, "*")

		redirect, err := params.Source.GetRedirectForKey(ctx, redirectKey)
		if err != nil {
			params.Logger.InfoContext(
				ctx,
				"Received unknown redirect key",
				slog.String("redirectKey", redirectKey),
			)

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Reason: "No matching redirect",
			})

			return
		}

		redirectURL, err := normalizeURL(redirect.URLPattern + "/" + redirectContextPath)
		if err != nil {
			params.Logger.ErrorContext(
				ctx,
				"Unable to sanitize redirect URL",
				logging.Err(err),
			)

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Reason: "Unable to get redirect, try again later",
			})

			return
		}

		w.Header().Set("Location", redirectURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func normalizeURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("unable to sanitize redirect URL: %w", err)
	}

	parsedURL.Path = strings.ReplaceAll(parsedURL.Path, "//", "/")

	return parsedURL.String(), nil
}
