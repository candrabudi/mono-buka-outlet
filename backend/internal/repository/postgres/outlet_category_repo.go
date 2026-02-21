package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type OutletCategoryRepo struct {
	db *sql.DB
}

func NewOutletCategoryRepo(db *sql.DB) *OutletCategoryRepo {
	return &OutletCategoryRepo{db: db}
}

func (r *OutletCategoryRepo) Create(ctx context.Context, cat *entity.OutletCategory) error {
	query := `INSERT INTO outlet_categories (id, name, slug, icon, description, is_active, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		cat.ID, cat.Name, cat.Slug, cat.Icon, cat.Description, cat.IsActive, cat.SortOrder,
		cat.CreatedAt, cat.UpdatedAt,
	)
	return err
}

func (r *OutletCategoryRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.OutletCategory, error) {
	query := `SELECT id, name, slug, icon, description, is_active, sort_order, created_at, updated_at
		FROM outlet_categories WHERE id = $1`
	cat := &entity.OutletCategory{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.Description,
		&cat.IsActive, &cat.SortOrder, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("kategori tidak ditemukan")
	}
	return cat, err
}

func (r *OutletCategoryRepo) FindBySlug(ctx context.Context, slug string) (*entity.OutletCategory, error) {
	query := `SELECT id, name, slug, icon, description, is_active, sort_order, created_at, updated_at
		FROM outlet_categories WHERE slug = $1`
	cat := &entity.OutletCategory{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.Description,
		&cat.IsActive, &cat.SortOrder, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("kategori tidak ditemukan")
	}
	return cat, err
}

func (r *OutletCategoryRepo) FindAll(ctx context.Context, activeOnly bool) ([]*entity.OutletCategory, error) {
	query := `SELECT id, name, slug, icon, description, is_active, sort_order, created_at, updated_at
		FROM outlet_categories`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY sort_order ASC, name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*entity.OutletCategory
	for rows.Next() {
		cat := &entity.OutletCategory{}
		if err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.Description,
			&cat.IsActive, &cat.SortOrder, &cat.CreatedAt, &cat.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

func (r *OutletCategoryRepo) Update(ctx context.Context, cat *entity.OutletCategory) error {
	query := `UPDATE outlet_categories SET name=$1, slug=$2, icon=$3, description=$4, is_active=$5, sort_order=$6, updated_at=$7
		WHERE id=$8`
	_, err := r.db.ExecContext(ctx, query,
		cat.Name, cat.Slug, cat.Icon, cat.Description, cat.IsActive, cat.SortOrder,
		time.Now(), cat.ID,
	)
	return err
}

func (r *OutletCategoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if any outlets use this category
	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outlets WHERE category_id = $1 AND deleted_at IS NULL", id).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("kategori masih digunakan oleh %d outlet", count)
	}
	_, err := r.db.ExecContext(ctx, "DELETE FROM outlet_categories WHERE id = $1", id)
	return err
}
