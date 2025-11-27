package service

import (
	"context"
	"ucode/ucode_go_sms_service/config"
	"ucode/ucode_go_sms_service/genproto/sms_service"
	"ucode/ucode_go_sms_service/pkg/logger"
	"ucode/ucode_go_sms_service/storage"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SendService struct...
type sendService struct {
	strg storage.StorageI
	cfg  config.Config
	log  logger.LoggerI
	sms_service.UnimplementedSmsServiceServer
}

// NewSendService ...
func NewSendService(cfg config.Config, log logger.LoggerI, strg storage.StorageI) *sendService {
	return &sendService{
		strg: strg,
		cfg:  cfg,
		log:  log,
	}
}

// Send ...
func (s *sendService) Send(ctx context.Context, req *sms_service.Sms) (*sms_service.GetSmsRequest, error) {
	s.log.Info("---Send--->", logger.Any("req", req))
	resp := &sms_service.GetSmsRequest{}

	resp, err := s.strg.Sms().Send(ctx, req)
	if err != nil {
		s.log.Error("!!!Send sms error--->", logger.Error(err))
		return resp, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *sendService) ConfirmOtp(ctx context.Context, req *sms_service.ConfirmOtpRequest) (resp *empty.Empty, err error) {
	s.log.Info("---ConfirmOtp--->", logger.Any("req", req))

	resp, err = s.strg.Sms().ConfirmOtp(ctx, req)
	if err != nil {
		s.log.Error("!!!ConfirmOtp--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return
}
