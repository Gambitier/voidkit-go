package server

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/Gambitier/voidkitgo/internal/config"
	"github.com/Gambitier/voidkitgo/internal/services"
	"github.com/sirupsen/logrus"
)

// Server represents the main server that coordinates HTTP and gRPC servers
type Server struct {
	config     *config.Config
	logger     *logrus.Logger
	httpServer *httpServer
	grpcServer GrpcServer
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, logger *logrus.Logger) *Server {
	return &Server{
		config: cfg,
		logger: logger,
	}
}

// recoverPanic recovers from panics and logs the error
func (s *Server) recoverPanic() {
	if r := recover(); r != nil {
		s.logger.Errorf("Recovered from panic in main thread: %v\nStack trace:\n%s", r, debug.Stack())
	}
}

// Start starts both HTTP and gRPC servers
func (s *Server) Start(ctx context.Context) error {
	// Add panic recovery for the main thread
	defer s.recoverPanic()

	// Create services
	services := services.NewServices()

	// Initialize servers
	s.httpServer = NewHTTPServer(HttpServerParams{
		services:  services,
		logger:    s.logger,
		serverEnv: s.config.Server.Env,
	})
	s.grpcServer = NewGrpcServer(GrpcServerParams{
		Services:  services,
		Logger:    s.logger,
		ServerEnv: s.config.Server.Env,
	})

	// Start servers in goroutines
	errChan := make(chan error, 2)

	// Start HTTP server with panic recovery
	go func() {
		defer s.recoverPanic()
		s.logger.Infof("Starting HTTP server on port %d", s.config.Server.HTTP.Port)
		if err := s.httpServer.Start(&s.config.Server.HTTP); err != nil {
			s.logger.Errorf("HTTP server error: %v", err)
			errChan <- err
		}
	}()

	// Start gRPC server with panic recovery
	go func() {
		defer s.recoverPanic()
		s.logger.Infof("Starting gRPC server on port %d", s.config.Server.GRPC.Port)
		if err := s.grpcServer.Start(s.config.Server.GRPC.Port); err != nil {
			s.logger.Errorf("gRPC server error: %v", err)
			errChan <- err
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received or a server error occurs
	select {
	case err := <-errChan:
		s.logger.Errorf("Server error: %v", err)
	case sig := <-sigChan:
		s.logger.Infof("Received signal: %v", sig)
	}

	// Graceful shutdown
	return s.Shutdown(ctx)
}

// Shutdown gracefully shuts down both servers
func (s *Server) Shutdown(ctx context.Context) error {
	defer s.recoverPanic()

	s.logger.Info("Shutting down servers...")
	// Create a timeout context for shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			s.logger.Errorf("HTTP server shutdown error: %v", err)
		}
	}

	// Shutdown gRPC server
	if s.grpcServer != nil {
		s.grpcServer.Shutdown()
	}

	s.logger.Info("Server shutdown complete")

	return nil
}
