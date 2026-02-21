package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type OTPRepo struct {
	db *sql.DB
}

func NewOTPRepo(db *sql.DB) *OTPRepo {
	return &OTPRepo{db: db}
}

func (r *OTPRepo) Create(ctx context.Context, otp *entity.OTPCode) error {
	query := `INSERT INTO otp_codes (id, user_id, email, code, purpose, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		otp.ID, otp.UserID, otp.Email, otp.Code, otp.Purpose, otp.ExpiresAt, otp.CreatedAt)
	return err
}

func (r *OTPRepo) FindLatestByEmailAndPurpose(ctx context.Context, email, purpose string) (*entity.OTPCode, error) {
	query := `SELECT id, user_id, email, code, purpose, expires_at, used_at, created_at
		FROM otp_codes WHERE email = $1 AND purpose = $2 AND used_at IS NULL
		ORDER BY created_at DESC LIMIT 1`
	otp := &entity.OTPCode{}
	err := r.db.QueryRowContext(ctx, query, email, purpose).Scan(
		&otp.ID, &otp.UserID, &otp.Email, &otp.Code, &otp.Purpose,
		&otp.ExpiresAt, &otp.UsedAt, &otp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *OTPRepo) FindByCode(ctx context.Context, email, code, purpose string) (*entity.OTPCode, error) {
	query := `SELECT id, user_id, email, code, purpose, expires_at, used_at, created_at
		FROM otp_codes WHERE email = $1 AND code = $2 AND purpose = $3 AND used_at IS NULL
		ORDER BY created_at DESC LIMIT 1`
	otp := &entity.OTPCode{}
	err := r.db.QueryRowContext(ctx, query, email, code, purpose).Scan(
		&otp.ID, &otp.UserID, &otp.Email, &otp.Code, &otp.Purpose,
		&otp.ExpiresAt, &otp.UsedAt, &otp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *OTPRepo) MarkUsed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE otp_codes SET used_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *OTPRepo) InvalidateAll(ctx context.Context, email, purpose string) error {
	query := `UPDATE otp_codes SET used_at = $1 WHERE email = $2 AND purpose = $3 AND used_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), email, purpose)
	return err
}
