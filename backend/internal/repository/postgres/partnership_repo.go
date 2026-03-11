package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type PartnershipRepo struct {
	db *sql.DB
}

func NewPartnershipRepo(db *sql.DB) *PartnershipRepo {
	return &PartnershipRepo{db: db}
}

func (r *PartnershipRepo) Create(ctx context.Context, p *entity.Partnership) error {
	query := `INSERT INTO partnerships (id, brand_id, mitra_id, affiliator_id, outlet_id, package_id, progress_percentage, status, start_date, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.BrandID, p.MitraID, p.AffiliatorID, p.OutletID, p.PackageID,
		p.ProgressPercentage, p.Status, p.StartDate, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PartnershipRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Partnership, error) {
	p := &entity.Partnership{}
	p.Mitra = &entity.User{}

	query := `SELECT p.id, p.brand_id, p.mitra_id, p.affiliator_id, p.outlet_id, p.package_id,
			  p.progress_percentage, p.status, p.start_date, p.created_at, p.updated_at,
			  u.id, u.name, u.email, u.phone,
			  au.id, au.name, au.email,
			  o.name,
			  op.name, op.price
			  FROM partnerships p
			  JOIN users u ON p.mitra_id = u.id
			  LEFT JOIN users au ON p.affiliator_id = au.id
			  LEFT JOIN outlets o ON p.outlet_id = o.id
			  LEFT JOIN outlet_packages op ON p.package_id = op.id
			  WHERE p.id = $1 AND p.deleted_at IS NULL`

	var affiliatorID sql.NullString
	var affiliatorName, affiliatorEmail sql.NullString
	var outletName, pkgName sql.NullString
	var pkgPrice sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.BrandID, &p.MitraID, &p.AffiliatorID, &p.OutletID, &p.PackageID,
		&p.ProgressPercentage, &p.Status, &p.StartDate, &p.CreatedAt, &p.UpdatedAt,
		&p.Mitra.ID, &p.Mitra.Name, &p.Mitra.Email, &p.Mitra.Phone,
		&affiliatorID, &affiliatorName, &affiliatorEmail,
		&outletName,
		&pkgName, &pkgPrice,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("partnership not found")
	}
	if err != nil {
		return nil, err
	}
	if affiliatorName.Valid {
		p.Affiliator = &entity.User{Name: affiliatorName.String, Email: affiliatorEmail.String}
		if affiliatorID.Valid {
			uid, _ := uuid.Parse(affiliatorID.String)
			p.Affiliator.ID = uid
		}
	}
	if outletName.Valid {
		p.Outlet = &entity.Outlet{Name: outletName.String}
	}
	if pkgName.Valid {
		p.Package = &entity.OutletPackage{Name: pkgName.String, Price: pkgPrice.Int64}
	}
	return p, nil
}

func (r *PartnershipRepo) FindAll(ctx context.Context, brandID *uuid.UUID, mitraID *uuid.UUID, page, limit int) ([]*entity.Partnership, int, error) {
	var total int
	countQuery := "SELECT COUNT(*) FROM partnerships WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if brandID != nil {
		countQuery += fmt.Sprintf(" AND brand_id = $%d", argIdx)
		args = append(args, *brandID)
		argIdx++
	}
	if mitraID != nil {
		countQuery += fmt.Sprintf(" AND mitra_id = $%d", argIdx)
		args = append(args, *mitraID)
		argIdx++
	}

	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT p.id, p.brand_id, p.mitra_id, p.affiliator_id, p.outlet_id, p.package_id,
			  p.progress_percentage, p.status, p.start_date, p.created_at, p.updated_at,
			  u.id, u.name, u.email,
			  au.id, au.name, au.email,
			  o.name, op.name, op.price
			  FROM partnerships p
			  JOIN users u ON p.mitra_id = u.id
			  LEFT JOIN users au ON p.affiliator_id = au.id
			  LEFT JOIN outlets o ON p.outlet_id = o.id
			  LEFT JOIN outlet_packages op ON p.package_id = op.id
			  WHERE p.deleted_at IS NULL`
	queryArgs := []interface{}{}
	queryArgIdx := 1

	if brandID != nil {
		query += fmt.Sprintf(" AND p.brand_id = $%d", queryArgIdx)
		queryArgs = append(queryArgs, *brandID)
		queryArgIdx++
	}
	if mitraID != nil {
		query += fmt.Sprintf(" AND p.mitra_id = $%d", queryArgIdx)
		queryArgs = append(queryArgs, *mitraID)
		queryArgIdx++
	}

	query += fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", queryArgIdx, queryArgIdx+1)
	queryArgs = append(queryArgs, limit, (page-1)*limit)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var partnerships []*entity.Partnership
	for rows.Next() {
		p := &entity.Partnership{}
		p.Mitra = &entity.User{}
		var affiliatorID sql.NullString
		var affiliatorName, affiliatorEmail sql.NullString
		var outletName, pkgName sql.NullString
		var pkgPrice sql.NullInt64
		err := rows.Scan(
			&p.ID, &p.BrandID, &p.MitraID, &p.AffiliatorID, &p.OutletID, &p.PackageID,
			&p.ProgressPercentage, &p.Status, &p.StartDate, &p.CreatedAt, &p.UpdatedAt,
			&p.Mitra.ID, &p.Mitra.Name, &p.Mitra.Email,
			&affiliatorID, &affiliatorName, &affiliatorEmail,
			&outletName, &pkgName, &pkgPrice,
		)
		if err != nil {
			return nil, 0, err
		}
		if affiliatorName.Valid {
			p.Affiliator = &entity.User{Name: affiliatorName.String, Email: affiliatorEmail.String}
			if affiliatorID.Valid {
				uid, _ := uuid.Parse(affiliatorID.String)
				p.Affiliator.ID = uid
			}
		}
		if outletName.Valid {
			p.Outlet = &entity.Outlet{Name: outletName.String}
		}
		if pkgName.Valid {
			p.Package = &entity.OutletPackage{Name: pkgName.String, Price: pkgPrice.Int64}
		}
		partnerships = append(partnerships, p)
	}

	return partnerships, total, nil
}

func (r *PartnershipRepo) FindByMitraID(ctx context.Context, mitraID uuid.UUID) ([]*entity.Partnership, error) {
	query := `SELECT p.id, p.brand_id, p.mitra_id, p.affiliator_id, p.outlet_id, p.package_id,
			  p.progress_percentage, p.status, p.start_date, p.created_at, p.updated_at,
			  au.id, au.name, au.email,
			  o.id, o.name, o.description,
			  op.id, op.name, op.price
			  FROM partnerships p
			  LEFT JOIN users au ON p.affiliator_id = au.id
			  LEFT JOIN outlets o ON p.outlet_id = o.id
			  LEFT JOIN outlet_packages op ON p.package_id = op.id
			  WHERE p.mitra_id = $1 AND p.deleted_at IS NULL
			  ORDER BY p.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, mitraID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partnerships []*entity.Partnership
	for rows.Next() {
		p := &entity.Partnership{}
		var affiliatorID, affiliatorName, affiliatorEmail sql.NullString
		var outletID, outletName, outletDesc sql.NullString
		var pkgID, pkgName sql.NullString
		var pkgPrice sql.NullInt64
		err := rows.Scan(
			&p.ID, &p.BrandID, &p.MitraID, &p.AffiliatorID, &p.OutletID, &p.PackageID,
			&p.ProgressPercentage, &p.Status, &p.StartDate, &p.CreatedAt, &p.UpdatedAt,
			&affiliatorID, &affiliatorName, &affiliatorEmail,
			&outletID, &outletName, &outletDesc,
			&pkgID, &pkgName, &pkgPrice,
		)
		if err != nil {
			return nil, err
		}
		if affiliatorName.Valid {
			p.Affiliator = &entity.User{Name: affiliatorName.String, Email: affiliatorEmail.String}
			if affiliatorID.Valid {
				uid, _ := uuid.Parse(affiliatorID.String)
				p.Affiliator.ID = uid
			}
		}
		if outletName.Valid {
			p.Outlet = &entity.Outlet{Name: outletName.String}
			if outletID.Valid {
				uid, _ := uuid.Parse(outletID.String)
				p.Outlet.ID = uid
			}
			if outletDesc.Valid {
				p.Outlet.Description = outletDesc.String
			}
		}
		if pkgName.Valid {
			p.Package = &entity.OutletPackage{Name: pkgName.String, Price: pkgPrice.Int64}
			if pkgID.Valid {
				uid, _ := uuid.Parse(pkgID.String)
				p.Package.ID = uid
			}
		}
		partnerships = append(partnerships, p)
	}

	return partnerships, nil
}

func (r *PartnershipRepo) Update(ctx context.Context, p *entity.Partnership) error {
	query := `UPDATE partnerships SET progress_percentage = $1, status = $2, start_date = $3, updated_at = $4 WHERE id = $5 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, p.ProgressPercentage, p.Status, p.StartDate, time.Now(), p.ID)
	return err
}

func (r *PartnershipRepo) UpdateProgress(ctx context.Context, id uuid.UUID, progress int, status string) error {
	query := `UPDATE partnerships SET progress_percentage = $1, status = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, progress, status, time.Now(), id)
	return err
}

func (r *PartnershipRepo) FindByAffiliatorID(ctx context.Context, affiliatorID uuid.UUID, page, limit int) ([]*entity.Partnership, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM partnerships WHERE affiliator_id = $1 AND deleted_at IS NULL`
	if err := r.db.QueryRowContext(ctx, countQuery, affiliatorID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT p.id, p.brand_id, p.mitra_id, p.affiliator_id, p.outlet_id, p.package_id,
			  p.progress_percentage, p.status, p.start_date, p.created_at, p.updated_at,
			  u.id, u.name, u.email,
			  o.name, op.name, op.price
			  FROM partnerships p
			  JOIN users u ON p.mitra_id = u.id
			  LEFT JOIN outlets o ON p.outlet_id = o.id
			  LEFT JOIN outlet_packages op ON p.package_id = op.id
			  WHERE p.affiliator_id = $1 AND p.deleted_at IS NULL
			  ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, affiliatorID, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var partnerships []*entity.Partnership
	for rows.Next() {
		p := &entity.Partnership{}
		p.Mitra = &entity.User{}
		var outletName, pkgName sql.NullString
		var pkgPrice sql.NullInt64
		err := rows.Scan(
			&p.ID, &p.BrandID, &p.MitraID, &p.AffiliatorID, &p.OutletID, &p.PackageID,
			&p.ProgressPercentage, &p.Status, &p.StartDate, &p.CreatedAt, &p.UpdatedAt,
			&p.Mitra.ID, &p.Mitra.Name, &p.Mitra.Email,
			&outletName, &pkgName, &pkgPrice,
		)
		if err != nil {
			return nil, 0, err
		}
		if outletName.Valid {
			p.Outlet = &entity.Outlet{Name: outletName.String}
		}
		if pkgName.Valid {
			p.Package = &entity.OutletPackage{Name: pkgName.String, Price: pkgPrice.Int64}
		}
		partnerships = append(partnerships, p)
	}

	return partnerships, total, nil
}
