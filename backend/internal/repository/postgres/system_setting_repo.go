package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/franchise-system/backend/internal/entity"
)

type SystemSettingRepo struct {
	db *sql.DB
}

func NewSystemSettingRepo(db *sql.DB) *SystemSettingRepo {
	return &SystemSettingRepo{db: db}
}

func (r *SystemSettingRepo) FindAll(ctx context.Context) ([]*entity.SystemSetting, error) {
	query := `SELECT id, key, value, group_name, label, description, created_at, updated_at 
			  FROM system_settings ORDER BY group_name, key`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []*entity.SystemSetting
	for rows.Next() {
		s := &entity.SystemSetting{}
		if err := rows.Scan(&s.ID, &s.Key, &s.Value, &s.GroupName, &s.Label, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

func (r *SystemSettingRepo) FindByGroup(ctx context.Context, group string) ([]*entity.SystemSetting, error) {
	query := `SELECT id, key, value, group_name, label, description, created_at, updated_at 
			  FROM system_settings WHERE group_name = $1 ORDER BY key`
	rows, err := r.db.QueryContext(ctx, query, group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []*entity.SystemSetting
	for rows.Next() {
		s := &entity.SystemSetting{}
		if err := rows.Scan(&s.ID, &s.Key, &s.Value, &s.GroupName, &s.Label, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

func (r *SystemSettingRepo) FindByKey(ctx context.Context, key string) (*entity.SystemSetting, error) {
	s := &entity.SystemSetting{}
	query := `SELECT id, key, value, group_name, label, description, created_at, updated_at 
			  FROM system_settings WHERE key = $1`
	err := r.db.QueryRowContext(ctx, query, key).Scan(&s.ID, &s.Key, &s.Value, &s.GroupName, &s.Label, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("setting not found: %s", key)
	}
	return s, err
}

func (r *SystemSettingRepo) Upsert(ctx context.Context, key, value string) error {
	query := `INSERT INTO system_settings (key, value, updated_at)
			  VALUES ($1, $2, $3)
			  ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = $3`
	_, err := r.db.ExecContext(ctx, query, key, value, time.Now())
	return err
}

func (r *SystemSettingRepo) BulkUpsert(ctx context.Context, settings map[string]string) error {
	if len(settings) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO system_settings (key, value, updated_at) VALUES ($1, $2, $3)
			  ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at`
	now := time.Now()

	for k, v := range settings {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, query, k, v, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}
