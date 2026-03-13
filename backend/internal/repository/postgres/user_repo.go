package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, name, email, password, phone, role, referral_code, referred_by, is_active, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Name, user.Email, user.Password, user.Phone,
		user.Role, nilIfEmpty(user.ReferralCode), user.ReferredBy, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user := &entity.User{}
	var referralCode sql.NullString
	var referredBy sql.NullString
	query := `SELECT id, name, email, password, phone, role, referral_code, referred_by, is_active, created_at, updated_at 
			  FROM users WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone,
		&user.Role, &referralCode, &referredBy, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	if referralCode.Valid {
		user.ReferralCode = referralCode.String
	}
	if referredBy.Valid {
		uid, _ := uuid.Parse(referredBy.String)
		user.ReferredBy = &uid
	}
	return user, err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	var referralCode sql.NullString
	var referredBy sql.NullString
	query := `SELECT id, name, email, password, phone, role, referral_code, referred_by, is_active, created_at, updated_at 
			  FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone,
		&user.Role, &referralCode, &referredBy, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	if referralCode.Valid {
		user.ReferralCode = referralCode.String
	}
	if referredBy.Valid {
		uid, _ := uuid.Parse(referredBy.String)
		user.ReferredBy = &uid
	}
	return user, err
}

func (r *UserRepo) FindAll(ctx context.Context, role string, page, limit int) ([]*entity.User, int, error) {
	var total int
	countQuery := "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if role != "" {
		countQuery += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, role)
		argIdx++
	}

	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT id, name, email, phone, role, referral_code, referred_by, is_active, created_at, updated_at FROM users WHERE deleted_at IS NULL"
	queryArgs := []interface{}{}
	queryArgIdx := 1

	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", queryArgIdx)
		queryArgs = append(queryArgs, role)
		queryArgIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", queryArgIdx, queryArgIdx+1)
	queryArgs = append(queryArgs, limit, (page-1)*limit)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		var referralCode sql.NullString
		var referredBy sql.NullString
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Phone,
			&user.Role, &referralCode, &referredBy, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		if referralCode.Valid {
			user.ReferralCode = referralCode.String
		}
		if referredBy.Valid {
			uid, _ := uuid.Parse(referredBy.String)
			user.ReferredBy = &uid
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET name = $1, email = $2, phone = $3, role = $4, is_active = $5, updated_at = $6, password = $7, referral_code = $8, referred_by = $9
			  WHERE id = $10 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query,
		user.Name, user.Email, user.Phone, user.Role, user.IsActive, time.Now(), user.Password, nilIfEmpty(user.ReferralCode), user.ReferredBy, user.ID,
	)
	return err
}

func (r *UserRepo) FindByReferralCode(ctx context.Context, code string) (*entity.User, error) {
	user := &entity.User{}
	var referralCode sql.NullString
	query := `SELECT id, name, email, password, phone, role, referral_code, is_active, created_at, updated_at 
			  FROM users WHERE referral_code = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone,
		&user.Role, &referralCode, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("affiliator tidak ditemukan")
	}
	if referralCode.Valid {
		user.ReferralCode = referralCode.String
	}
	return user, err
}

func (r *UserRepo) CountByRole(ctx context.Context, role string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE role = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, role).Scan(&count)
	return count, err
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE users SET deleted_at = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
