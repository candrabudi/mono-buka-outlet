package postgres

import (
	"context"
	"database/sql"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type ActivityLogRepo struct {
	db *sql.DB
}

func NewActivityLogRepo(db *sql.DB) *ActivityLogRepo {
	return &ActivityLogRepo{db: db}
}

func (r *ActivityLogRepo) Create(ctx context.Context, l *entity.ActivityLog) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO activity_logs (id, entity_type, entity_id, action, description, old_value, new_value, performed_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, l.ID, l.EntityType, l.EntityID, l.Action, l.Description, l.OldValue, l.NewValue, l.PerformedBy, l.CreatedAt)
	return err
}

func (r *ActivityLogRepo) FindByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.ActivityLog, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT al.id, al.entity_type, al.entity_id, al.action, al.description, al.old_value, al.new_value, al.performed_by, al.created_at, u.name FROM activity_logs al JOIN users u ON al.performed_by = u.id WHERE al.entity_type=$1 AND al.entity_id=$2 ORDER BY al.created_at DESC`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entity.ActivityLog
	for rows.Next() {
		l := &entity.ActivityLog{Performer: &entity.User{}}
		if err := rows.Scan(&l.ID, &l.EntityType, &l.EntityID, &l.Action, &l.Description, &l.OldValue, &l.NewValue, &l.PerformedBy, &l.CreatedAt, &l.Performer.Name); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}
