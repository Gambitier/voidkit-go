package server

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/Gambitier/voidkitgo/internal/config"
	"github.com/Gambitier/voidkitgo/internal/services"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GrpcServer interface {
	// Start starts the server on the specified port
	Start(port int) error
	// Shutdown gracefully stops the server
	Shutdown()
	// Port returns the port the server is listening on
	Port() int
}

type grpcServer struct {
	server    *grpc.Server
	serverEnv config.Environment
	port      int
	logger    *logrus.Logger
}

type GrpcServerParams struct {
	Services  *services.Services
	Logger    *logrus.Logger
	ServerEnv config.Environment
}

// panicRecoveryUnaryInterceptor returns a new unary server interceptor for panic recovery
func panicRecoveryUnaryInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("Recovered from panic in gRPC handler: %v\nStack trace:\n%s", r, debug.Stack())
				err = status.Errorf(codes.Internal, "Internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

// panicRecoveryStreamInterceptor returns a new stream server interceptor for panic recovery
func panicRecoveryStreamInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("Recovered from panic in gRPC stream handler: %v\nStack trace:\n%s", r, debug.Stack())
				err = status.Errorf(codes.Internal, "Internal server error")
			}
		}()
		return handler(srv, stream)
	}
}

// NewGrpcServer creates a new gRPC server
func NewGrpcServer(params GrpcServerParams) GrpcServer {
	return &grpcServer{
		logger:    params.Logger,
		serverEnv: params.ServerEnv,
	}
}

// Start starts the gRPC server
func (s *grpcServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Create gRPC server with panic recovery interceptors
	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(panicRecoveryUnaryInterceptor(s.logger)),
		grpc.StreamInterceptor(panicRecoveryStreamInterceptor(s.logger)),
	)

	if s.serverEnv.IsDevelopment() {
		reflection.Register(s.server)
	}

	s.port = port
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the gRPC server
func (s *grpcServer) Shutdown() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

// Port returns the port the server is listening on
func (s *grpcServer) Port() int {
	return s.port
}
