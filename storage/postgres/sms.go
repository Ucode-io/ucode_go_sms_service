package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ucode/ucode_go_sms_service/config"
	"ucode/ucode_go_sms_service/genproto/sms_service"
	"ucode/ucode_go_sms_service/storage"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type smsRepo struct {
	db *pgxpool.Pool
}

// NewSmsRepo ...
func NewSmsRepo(db *pgxpool.Pool) storage.SmsRepoI {
	return &smsRepo{db: db}
}

// GetNotSent ...
func (cm *smsRepo) GetNotSent(ctx context.Context) ([]*sms_service.Sms, error) {
	var smss []*sms_service.Sms

	query := `SELECT  
			id,
			recipient,
			text,
			send_count,
			otp,
			expires_at,
			type,
			dev_email,
			dev_email_password,
			originator
	FROM sms_send
	WHERE sent_at IS NULL and send_count < 4`
	rows, err := cm.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			sms                                   sms_service.Sms
			expiresAt                             sql.NullString
			typeOfSms, devEmailPassword, devEmail sql.NullString
		)
		if err := rows.Scan(
			&sms.Id,
			&sms.Recipient,
			&sms.Text,
			&sms.SendCount,
			&sms.Otp,
			&expiresAt,
			&typeOfSms,
			&devEmail,
			&devEmailPassword,
			&sms.Originator,
		); err != nil {
			return nil, err
		}

		if expiresAt.Valid {
			sms.ExpiresAt = expiresAt.String
		}
		if typeOfSms.Valid {
			sms.Type = typeOfSms.String
		}
		if devEmail.Valid {
			sms.DevEmail = devEmail.String
		}
		if devEmailPassword.Valid {
			sms.DevEmailPassword = devEmailPassword.String
		}

		smss = append(smss, &sms)
	}
	return smss, err
}

// MakeSent ...
func (cm *smsRepo) MakeSent(ctx context.Context, ID string) error {
	makesent := `UPDATE sms_send SET sent_at = CURRENT_TIMESTAMP where id = $1`
	_, err := cm.db.Exec(ctx, makesent, ID)

	return err
}

func (cm *smsRepo) IncrementSendCount(ctx context.Context, ID string) error {
	query := `UPDATE sms_send SET send_count = send_count + 1 where id = $1`
	_, err := cm.db.Exec(ctx, query, ID)
	return err
}

func (cm *smsRepo) Send(ctx context.Context, req *sms_service.Sms) (*sms_service.GetSmsRequest, error) {
	resp := &sms_service.GetSmsRequest{}
	sendID, err := uuid.NewRandom()
	if err != nil {
		return resp, err
	}

	query := fmt.Sprintf(`INSERT INTO
		sms_send
		(
			id,
			text,
			recipient,
			otp,
			expires_at,
			type,
			dev_email,
			dev_email_password,
			originator
		)
		values($1, $2, $3, $4, now() + interval '%d minute', $5, $6, $7, $8)`, config.ExpireTimeout)

	_, err = cm.db.Exec(ctx, query, sendID, req.Text, req.Recipient, req.Otp, req.Type, req.DevEmail, req.DevEmailPassword, req.Originator)
	if err != nil {
		return resp, err
	}

	resp.SmsId = sendID.String()

	return resp, nil
}

func (cm *smsRepo) ConfirmOtp(ctx context.Context, req *sms_service.ConfirmOtpRequest) (resp *empty.Empty, err error) {
	var (
		exist bool
	)

	resp = &emptypb.Empty{}

	query := `
		SELECT EXISTS(SELECT 1 FROM sms_send WHERE id = $1 AND otp = $2 AND expires_at >= CURRENT_TIMESTAMP)
	`

	err = cm.db.QueryRow(ctx, query, req.SmsId, req.Otp).Scan(&exist)
	if err != nil {
		return resp, errors.Wrap(err, "Error while scanning exists")
	}

	if !exist {
		return resp, errors.New("wrong otp or expired")
	}

	return
}
