package http

import (
	"net/http"

	"github.com/Gambitier/voidkitgo/internal/server/handlers/http/health"
	"github.com/gorilla/mux"
)

type HttpHandlers struct {
	healthHandler *health.HealthHandler
}

func NewHttpHandlers(server *http.Server) *HttpHandlers {
	return &HttpHandlers{
		healthHandler: health.NewHealthHandler(server),
	}
}

func (h *HttpHandlers) RegisterRoutes(router *mux.Router) {
	h.healthHandler.RegisterRoutes(router)
}
