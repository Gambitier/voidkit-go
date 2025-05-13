package health

import (
	"encoding/json"
	"net/http"

	"github.com/Gambitier/voidkitgo/internal/server/handlers/http/common"
	"github.com/gorilla/mux"
)

type healthHandler struct {
	server *http.Server
}

func NewHealthHandler(server *http.Server) common.HttpHandler {
	return &healthHandler{server: server}
}

// register routes
func (h *healthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.HandleHealthCheck).Methods(http.MethodGet)
}

func (h *healthHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
