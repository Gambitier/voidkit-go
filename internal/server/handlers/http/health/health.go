package health

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type HealthHandler struct {
	server *http.Server
}

func NewHealthHandler(server *http.Server) *HealthHandler {
	return &HealthHandler{server: server}
}

// register routes
func (h *HealthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.HandleHealthCheck).Methods(http.MethodGet)
}

func (h *HealthHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
