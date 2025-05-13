package common

import (
	"context"

	"github.com/Gambitier/voidkitgo/internal/services"
	"github.com/Gambitier/voidkitgo/pkg/proto/common"
)

type Handler struct {
	common.UnimplementedCommonServiceServer
	services *services.Services
}

// NewCommonServiceHandler creates a new common service handler
func NewCommonServiceHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) HealthCheck(
	ctx context.Context,
	req *common.HealthCheckRequest,
) (*common.HealthCheckResponse, error) {
	return &common.HealthCheckResponse{Status: true}, nil
}
