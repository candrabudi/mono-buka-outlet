package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type OutletRepo struct {
	db *sql.DB
}

func NewOutletRepo(db *sql.DB) *OutletRepo {
	return &OutletRepo{db: db}
}

// Select columns include a LEFT JOIN on outlet_categories to populate category_name and category_icon
var outletSelectQuery = `SELECT o.id, o.name, o.slug, COALESCE(o.logo,''), COALESCE(o.banner,''), o.category_id, COALESCE(o.category,''),
	COALESCE(oc.name, o.category, ''), COALESCE(oc.icon, ''),
	COALESCE(o.description,''), COALESCE(o.short_description,''),
	o.minimum_investment, o.maximum_investment, o.profit_sharing_percentage, COALESCE(o.estimated_roi,''),
	COALESCE(o.location_requirement,''), COALESCE(o.address,''), COALESCE(o.city,''), COALESCE(o.province,''), o.latitude, o.longitude,
	COALESCE(o.contact_phone,''), COALESCE(o.contact_email,''), COALESCE(o.contact_whatsapp,''), COALESCE(o.website,''),
	o.is_active, o.is_featured, o.total_outlets, o.year_established, o.created_by,
	o.created_at, o.updated_at
FROM outlets o
LEFT JOIN outlet_categories oc ON o.category_id = oc.id`

func scanOutlet(row interface {
	Scan(dest ...interface{}) error
}) (*entity.Outlet, error) {
	o := &entity.Outlet{}
	err := row.Scan(
		&o.ID, &o.Name, &o.Slug, &o.Logo, &o.Banner, &o.CategoryID, &o.Category,
		&o.CategoryName, &o.CategoryIcon,
		&o.Description, &o.ShortDescription,
		&o.MinimumInvestment, &o.MaximumInvestment, &o.ProfitSharingPercentage, &o.EstimatedROI,
		&o.LocationRequirement, &o.Address, &o.City, &o.Province,
		&o.Latitude, &o.Longitude,
		&o.ContactPhone, &o.ContactEmail, &o.ContactWhatsapp, &o.Website,
		&o.IsActive, &o.IsFeatured, &o.TotalOutlets, &o.YearEstablished, &o.CreatedBy,
		&o.CreatedAt, &o.UpdatedAt,
	)
	return o, err
}

func (r *OutletRepo) Create(ctx context.Context, outlet *entity.Outlet) error {
	query := `INSERT INTO outlets (id, name, slug, logo, banner, category_id, category, description, short_description,
		minimum_investment, maximum_investment, profit_sharing_percentage, estimated_roi,
		location_requirement, address, city, province, latitude, longitude,
		contact_phone, contact_email, contact_whatsapp, website,
		is_active, is_featured, total_outlets, year_established, created_by,
		created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30)`
	_, err := r.db.ExecContext(ctx, query,
		outlet.ID, outlet.Name, outlet.Slug, outlet.Logo, outlet.Banner, outlet.CategoryID, outlet.Category,
		outlet.Description, outlet.ShortDescription,
		outlet.MinimumInvestment, outlet.MaximumInvestment, outlet.ProfitSharingPercentage, outlet.EstimatedROI,
		outlet.LocationRequirement, outlet.Address, outlet.City, outlet.Province,
		outlet.Latitude, outlet.Longitude,
		outlet.ContactPhone, outlet.ContactEmail, outlet.ContactWhatsapp, outlet.Website,
		outlet.IsActive, outlet.IsFeatured, outlet.TotalOutlets, outlet.YearEstablished, outlet.CreatedBy,
		outlet.CreatedAt, outlet.UpdatedAt,
	)
	return err
}

func (r *OutletRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Outlet, error) {
	query := outletSelectQuery + ` WHERE o.id = $1 AND o.deleted_at IS NULL`
	o, err := scanOutlet(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("outlet not found")
	}
	return o, err
}

func (r *OutletRepo) FindBySlug(ctx context.Context, slug string) (*entity.Outlet, error) {
	query := outletSelectQuery + ` WHERE o.slug = $1 AND o.deleted_at IS NULL`
	o, err := scanOutlet(r.db.QueryRowContext(ctx, query, slug))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("outlet not found")
	}
	return o, err
}

func (r *OutletRepo) FindAll(ctx context.Context, filter repository.OutletFilter) ([]*entity.Outlet, int, error) {
	where := []string{"o.deleted_at IS NULL"}
	args := []interface{}{}
	argIdx := 1

	if filter.Search != "" {
		where = append(where, fmt.Sprintf("(LOWER(o.name) LIKE $%d OR LOWER(o.city) LIKE $%d OR LOWER(o.short_description) LIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+strings.ToLower(filter.Search)+"%")
		argIdx++
	}
	if filter.CategoryID != nil {
		where = append(where, fmt.Sprintf("o.category_id = $%d", argIdx))
		args = append(args, *filter.CategoryID)
		argIdx++
	}
	if filter.City != "" {
		where = append(where, fmt.Sprintf("LOWER(o.city) = $%d", argIdx))
		args = append(args, strings.ToLower(filter.City))
		argIdx++
	}
	if filter.Province != "" {
		where = append(where, fmt.Sprintf("LOWER(o.province) = $%d", argIdx))
		args = append(args, strings.ToLower(filter.Province))
		argIdx++
	}
	if filter.Active != nil {
		where = append(where, fmt.Sprintf("o.is_active = $%d", argIdx))
		args = append(args, *filter.Active)
		argIdx++
	}
	if filter.Featured != nil {
		where = append(where, fmt.Sprintf("o.is_featured = $%d", argIdx))
		args = append(args, *filter.Featured)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM outlets o WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page := filter.Page
	limit := filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	query := fmt.Sprintf(`%s WHERE %s ORDER BY o.is_featured DESC, o.created_at DESC LIMIT $%d OFFSET $%d`,
		outletSelectQuery, whereClause, argIdx, argIdx+1)
	queryArgs := append(args, limit, (page-1)*limit)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var outlets []*entity.Outlet
	for rows.Next() {
		o, err := scanOutlet(rows)
		if err != nil {
			return nil, 0, err
		}
		outlets = append(outlets, o)
	}

	return outlets, total, nil
}

func (r *OutletRepo) Update(ctx context.Context, outlet *entity.Outlet) error {
	query := `UPDATE outlets SET
		name=$1, slug=$2, logo=$3, banner=$4, category_id=$5, category=$6,
		description=$7, short_description=$8,
		minimum_investment=$9, maximum_investment=$10, profit_sharing_percentage=$11, estimated_roi=$12,
		location_requirement=$13, address=$14, city=$15, province=$16, latitude=$17, longitude=$18,
		contact_phone=$19, contact_email=$20, contact_whatsapp=$21, website=$22,
		is_active=$23, is_featured=$24, total_outlets=$25, year_established=$26,
		updated_at=$27
		WHERE id=$28 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query,
		outlet.Name, outlet.Slug, outlet.Logo, outlet.Banner, outlet.CategoryID, outlet.Category,
		outlet.Description, outlet.ShortDescription,
		outlet.MinimumInvestment, outlet.MaximumInvestment, outlet.ProfitSharingPercentage, outlet.EstimatedROI,
		outlet.LocationRequirement, outlet.Address, outlet.City, outlet.Province,
		outlet.Latitude, outlet.Longitude,
		outlet.ContactPhone, outlet.ContactEmail, outlet.ContactWhatsapp, outlet.Website,
		outlet.IsActive, outlet.IsFeatured, outlet.TotalOutlets, outlet.YearEstablished,
		time.Now(), outlet.ID,
	)
	return err
}

func (r *OutletRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE outlets SET deleted_at = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
