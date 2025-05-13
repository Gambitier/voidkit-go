package grpc

import (
	"github.com/Gambitier/voidkitgo/internal/server/handlers/grpc/common"
	"github.com/Gambitier/voidkitgo/internal/services"
	commonProto "github.com/Gambitier/voidkitgo/pkg/proto/common"
	"google.golang.org/grpc"
)

type GrpcHandlers struct {
	CommonServiceHandler *common.Handler
}

func NewGrpcHandlers(services *services.Services) *GrpcHandlers {
	commonServiceHandler := common.NewCommonServiceHandler(services)

	return &GrpcHandlers{
		CommonServiceHandler: commonServiceHandler,
	}
}

func (h *GrpcHandlers) RegisterServices(server *grpc.Server) {
	// register new service servers here
	commonProto.RegisterCommonServiceServer(server, h.CommonServiceHandler)
}
