package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type PartnershipApplicationRepo struct {
	db *sql.DB
}

func NewPartnershipApplicationRepo(db *sql.DB) *PartnershipApplicationRepo {
	return &PartnershipApplicationRepo{db: db}
}

func (r *PartnershipApplicationRepo) Create(ctx context.Context, app *entity.PartnershipApplication) error {
	query := `
		INSERT INTO partnership_applications (id, mitra_id, outlet_id, package_id, motivation, experience, proposed_location, investment_budget, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		app.ID, app.MitraID, app.OutletID, app.PackageID,
		app.Motivation, app.Experience, app.ProposedLocation, app.InvestmentBudget,
		app.Status, app.CreatedAt, app.UpdatedAt,
	)
	return err
}

func (r *PartnershipApplicationRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.PartnershipApplication, error) {
	query := `
		SELECT a.id, a.mitra_id, a.outlet_id, a.package_id,
		       a.motivation, a.experience, a.proposed_location, a.investment_budget,
		       a.status, a.admin_notes, a.reviewed_by, a.reviewed_at,
		       a.created_at, a.updated_at,
		       u.id, u.name, u.email, u.phone,
		       o.id, o.name, o.slug, o.logo, o.minimum_investment,
		       p.id, p.name, p.price
		FROM partnership_applications a
		JOIN users u ON u.id = a.mitra_id
		JOIN outlets o ON o.id = a.outlet_id
		JOIN outlet_packages p ON p.id = a.package_id
		WHERE a.id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var app entity.PartnershipApplication
	app.Mitra = &entity.User{}
	app.Outlet = &entity.Outlet{}
	app.Package = &entity.OutletPackage{}
	var adminNotes, reviewedBy sql.NullString
	var reviewedAt sql.NullTime

	err := row.Scan(
		&app.ID, &app.MitraID, &app.OutletID, &app.PackageID,
		&app.Motivation, &app.Experience, &app.ProposedLocation, &app.InvestmentBudget,
		&app.Status, &adminNotes, &reviewedBy, &reviewedAt,
		&app.CreatedAt, &app.UpdatedAt,
		&app.Mitra.ID, &app.Mitra.Name, &app.Mitra.Email, &app.Mitra.Phone,
		&app.Outlet.ID, &app.Outlet.Name, &app.Outlet.Slug, &app.Outlet.Logo, &app.Outlet.MinimumInvestment,
		&app.Package.ID, &app.Package.Name, &app.Package.Price,
	)
	if err != nil {
		return nil, err
	}

	if adminNotes.Valid {
		app.AdminNotes = adminNotes.String
	}
	if reviewedBy.Valid {
		uid, _ := uuid.Parse(reviewedBy.String)
		app.ReviewedBy = &uid
	}
	if reviewedAt.Valid {
		app.ReviewedAt = &reviewedAt.Time
	}

	return &app, nil
}

func (r *PartnershipApplicationRepo) FindByMitraID(ctx context.Context, mitraID uuid.UUID) ([]*entity.PartnershipApplication, error) {
	query := `
		SELECT a.id, a.mitra_id, a.outlet_id, a.package_id,
		       a.motivation, a.experience, a.proposed_location, a.investment_budget,
		       a.status, a.admin_notes, a.reviewed_at,
		       a.created_at, a.updated_at,
		       o.id, o.name, o.slug, o.logo, o.minimum_investment,
		       p.id, p.name, p.price
		FROM partnership_applications a
		JOIN outlets o ON o.id = a.outlet_id
		JOIN outlet_packages p ON p.id = a.package_id
		WHERE a.mitra_id = $1
		ORDER BY a.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, mitraID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.PartnershipApplication
	for rows.Next() {
		app := &entity.PartnershipApplication{}
		app.Outlet = &entity.Outlet{}
		app.Package = &entity.OutletPackage{}
		var adminNotes sql.NullString
		var reviewedAt sql.NullTime

		err := rows.Scan(
			&app.ID, &app.MitraID, &app.OutletID, &app.PackageID,
			&app.Motivation, &app.Experience, &app.ProposedLocation, &app.InvestmentBudget,
			&app.Status, &adminNotes, &reviewedAt,
			&app.CreatedAt, &app.UpdatedAt,
			&app.Outlet.ID, &app.Outlet.Name, &app.Outlet.Slug, &app.Outlet.Logo, &app.Outlet.MinimumInvestment,
			&app.Package.ID, &app.Package.Name, &app.Package.Price,
		)
		if err != nil {
			return nil, err
		}
		if adminNotes.Valid {
			app.AdminNotes = adminNotes.String
		}
		if reviewedAt.Valid {
			app.ReviewedAt = &reviewedAt.Time
		}
		results = append(results, app)
	}
	return results, nil
}

func (r *PartnershipApplicationRepo) FindAll(ctx context.Context, status string, page, limit int) ([]*entity.PartnershipApplication, int, error) {
	countQuery := `SELECT COUNT(*) FROM partnership_applications WHERE 1=1`
	dataQuery := `
		SELECT a.id, a.mitra_id, a.outlet_id, a.package_id,
		       a.motivation, a.experience, a.proposed_location, a.investment_budget,
		       a.status, a.admin_notes, a.reviewed_at,
		       a.created_at, a.updated_at,
		       u.id, u.name, u.email, u.phone,
		       o.id, o.name, o.slug, o.logo,
		       p.id, p.name, p.price
		FROM partnership_applications a
		JOIN users u ON u.id = a.mitra_id
		JOIN outlets o ON o.id = a.outlet_id
		JOIN outlet_packages p ON p.id = a.package_id
		WHERE 1=1
	`

	args := []interface{}{}
	argIdx := 1

	if status != "" {
		filter := fmt.Sprintf(" AND a.status = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		dataQuery += filter
		args = append(args, status)
		argIdx++
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery += fmt.Sprintf(" ORDER BY a.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*entity.PartnershipApplication
	for rows.Next() {
		app := &entity.PartnershipApplication{}
		app.Mitra = &entity.User{}
		app.Outlet = &entity.Outlet{}
		app.Package = &entity.OutletPackage{}
		var adminNotes sql.NullString
		var reviewedAt sql.NullTime

		err := rows.Scan(
			&app.ID, &app.MitraID, &app.OutletID, &app.PackageID,
			&app.Motivation, &app.Experience, &app.ProposedLocation, &app.InvestmentBudget,
			&app.Status, &adminNotes, &reviewedAt,
			&app.CreatedAt, &app.UpdatedAt,
			&app.Mitra.ID, &app.Mitra.Name, &app.Mitra.Email, &app.Mitra.Phone,
			&app.Outlet.ID, &app.Outlet.Name, &app.Outlet.Slug, &app.Outlet.Logo,
			&app.Package.ID, &app.Package.Name, &app.Package.Price,
		)
		if err != nil {
			return nil, 0, err
		}
		if adminNotes.Valid {
			app.AdminNotes = adminNotes.String
		}
		if reviewedAt.Valid {
			app.ReviewedAt = &reviewedAt.Time
		}
		results = append(results, app)
	}
	return results, total, nil
}

func (r *PartnershipApplicationRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status, adminNotes string, reviewedBy uuid.UUID) error {
	query := `
		UPDATE partnership_applications
		SET status = $1, admin_notes = $2, reviewed_by = $3, reviewed_at = $4, updated_at = $5
		WHERE id = $6
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, status, adminNotes, reviewedBy, now, now, id)
	return err
}
