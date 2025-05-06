package web

import (
	"log/slog"
	"net/http"

	"github.com/durandj/ek/internal/logging"
	"github.com/durandj/ek/internal/sources"
	"github.com/go-chi/render"
)

type RedirectsListEndpointParams struct {
	Logger *slog.Logger
	Source sources.Source
}

func RedirectsListEndpoint(params RedirectsListEndpointParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		redirects, err := params.Source.GetAllRedirects(ctx)
		if err != nil {
			params.Logger.ErrorContext(
				ctx,
				"Unable to retrieve all redirects",
				logging.Err(err),
			)

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Reason: "Internal server error, try again",
			})

			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, RedirectsListResponse{
			Redirects: redirects,
		})
	}
}

type RedirectsListResponse struct {
	Redirects map[string]sources.Redirect `json:"redirects"`
}

var _ render.Renderer = (*RedirectsListResponse)(nil)

func (response RedirectsListResponse) Render(
	w http.ResponseWriter,
	r *http.Request,
) error {
	return nil
}
