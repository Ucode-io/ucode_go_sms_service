package grpc

import (
	"ucode/ucode_go_sms_service/config"
	"ucode/ucode_go_sms_service/genproto/sms_service"
	"ucode/ucode_go_sms_service/grpc/service"
	"ucode/ucode_go_sms_service/pkg/logger"
	"ucode/ucode_go_sms_service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer()

	sms_service.RegisterSmsServiceServer(grpcServer, service.NewSendService(cfg, log, strg))

	reflection.Register(grpcServer)
	return
}
