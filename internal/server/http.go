package server

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Gambitier/voidkitgo/internal/config"
	httpHandlers "github.com/Gambitier/voidkitgo/internal/server/handlers/http"
	"github.com/Gambitier/voidkitgo/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// httpServer represents the HTTP server
type httpServer struct {
	router *mux.Router
	server *http.Server
	logger *logrus.Logger
}

type HttpServerParams struct {
	services  *services.Services
	logger    *logrus.Logger
	serverEnv config.Environment
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(params HttpServerParams) *httpServer {
	router := mux.NewRouter()
	server := &httpServer{
		router: router,
		logger: params.logger,
	}
	return server
}

// panicRecoveryMiddleware recovers from panics and logs the error
func (s *httpServer) panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Errorf("Recovered from panic in HTTP handler: %v\nStack trace:\n%s", err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Start starts the HTTP server
func (s *httpServer) Start(config *config.HTTPConfig) error {
	// Add panic recovery middleware
	handler := s.panicRecoveryMiddleware(s.router)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	httpHandlers := httpHandlers.NewHttpHandlers(s.server)
	httpHandlers.RegisterRoutes(s.router)

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *httpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
