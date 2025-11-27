package storage

import (
	"context"
	"ucode/ucode_go_sms_service/genproto/sms_service"

	"github.com/golang/protobuf/ptypes/empty"
)

type StorageI interface {
	CloseDB()
	Sms() SmsRepoI
}

// SmstorageI ...
type SmsRepoI interface {
	GetNotSent(ctx context.Context) ([]*sms_service.Sms, error)
	MakeSent(ctx context.Context, ID string) error
	IncrementSendCount(ctx context.Context, ID string) error
	Send(ctx context.Context, req *sms_service.Sms) (*sms_service.GetSmsRequest, error)
	ConfirmOtp(ctx context.Context, req *sms_service.ConfirmOtpRequest) (resp *empty.Empty, err error)
}
