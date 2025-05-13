package http

import (
	"net/http"

	"github.com/Gambitier/voidkitgo/internal/server/handlers/http/common"
	"github.com/Gambitier/voidkitgo/internal/server/handlers/http/health"
	"github.com/gorilla/mux"
)

type httpHandlers struct {
	healthHandler common.HttpHandler
}

func NewHttpHandlers(server *http.Server) common.HttpHandler {
	return &httpHandlers{
		healthHandler: health.NewHealthHandler(server),
	}
}

func (h *httpHandlers) RegisterRoutes(router *mux.Router) {
	h.healthHandler.RegisterRoutes(router)
}
