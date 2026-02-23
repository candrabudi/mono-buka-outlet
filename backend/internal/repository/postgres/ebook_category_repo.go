package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type EbookCategoryRepo struct {
	db *sql.DB
}

func NewEbookCategoryRepo(db *sql.DB) *EbookCategoryRepo {
	return &EbookCategoryRepo{db: db}
}

func (r *EbookCategoryRepo) Create(ctx context.Context, cat *entity.EbookCategory) error {
	if cat.ID == uuid.Nil {
		cat.ID = uuid.New()
	}
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()

	query := `INSERT INTO ebook_categories (id, name, slug, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		cat.ID, cat.Name, cat.Slug, cat.IsActive, cat.CreatedAt, cat.UpdatedAt)
	return err
}

func (r *EbookCategoryRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.EbookCategory, error) {
	cat := &entity.EbookCategory{}
	query := `SELECT id, name, slug, is_active, created_at, updated_at FROM ebook_categories WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("kategori ebook tidak ditemukan")
	}
	return cat, err
}

func (r *EbookCategoryRepo) FindBySlug(ctx context.Context, slug string) (*entity.EbookCategory, error) {
	cat := &entity.EbookCategory{}
	query := `SELECT id, name, slug, is_active, created_at, updated_at FROM ebook_categories WHERE slug = $1`
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("kategori ebook tidak ditemukan")
	}
	return cat, err
}

func (r *EbookCategoryRepo) FindAll(ctx context.Context, activeOnly bool) ([]*entity.EbookCategory, error) {
	query := `SELECT id, name, slug, is_active, created_at, updated_at FROM ebook_categories`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*entity.EbookCategory
	for rows.Next() {
		c := &entity.EbookCategory{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *EbookCategoryRepo) Update(ctx context.Context, cat *entity.EbookCategory) error {
	cat.UpdatedAt = time.Now()
	query := `UPDATE ebook_categories SET name=$1, slug=$2, is_active=$3, updated_at=$4 WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query, cat.Name, cat.Slug, cat.IsActive, cat.UpdatedAt, cat.ID)
	return err
}

func (r *EbookCategoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if used by any ebook
	var count int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM ebook_category_mapping WHERE category_id = $1", id).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("kategori masih digunakan oleh %d ebook", count)
	}
	_, err := r.db.ExecContext(ctx, "DELETE FROM ebook_categories WHERE id = $1", id)
	return err
}
