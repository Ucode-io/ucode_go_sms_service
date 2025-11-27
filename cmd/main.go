package main

import (
	"context"
	"net"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"ucode/ucode_go_sms_service/config"
	"ucode/ucode_go_sms_service/grpc"
	"ucode/ucode_go_sms_service/pkg/logger"
	"ucode/ucode_go_sms_service/pkg/sms"
	"ucode/ucode_go_sms_service/storage/postgres"
)

func main() {
	cfg := config.Load()

	loggerLevel := logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer logger.Cleanup(log)

	pgStore, err := postgres.NewPostgres(context.Background(), cfg)
	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}
	defer pgStore.CloseDB()

	grpcServer := grpc.SetUpServer(cfg, log, pgStore)

	smsDaemon := sms.Daemon{Conf: cfg, Strg: pgStore}
	go smsDaemon.Init()
	log.Info("Configs", logger.Any("config", cfg))
	lis, err := net.Listen("tcp", cfg.SmsGRPCPort)
	if err != nil {
		log.Panic("net.Listen", logger.Error(err))
	}

	log.Info("GRPC: Server being started...", logger.String("port", cfg.SmsGRPCPort))

	if err := grpcServer.Serve(lis); err != nil {
		log.Panic("grpcServer.Serve", logger.Error(err))

	}
}
