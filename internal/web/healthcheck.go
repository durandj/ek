package web

import (
	"net/http"

	"github.com/go-chi/render"
)

func HealthcheckEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, HealthcheckResponse{
			Status: "OK",
		})
	}
}

type HealthcheckResponse struct {
	Status string `json:"status"`
}
