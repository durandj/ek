package web

import (
	"log/slog"
	"net/http"
	"strings"

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

		redirectURL := normalizeURL(redirect.URLPattern + "/" + redirectContextPath)

		w.Header().Set("Location", redirectURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func normalizeURL(url string) string {
	return strings.ReplaceAll(url, "//", "/")
}
