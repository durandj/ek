package web

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

var _ render.Renderer = (*ErrorResponse)(nil)

func (response ErrorResponse) Render(
	w http.ResponseWriter,
	r *http.Request,
) error {
	return nil
}
