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
	query := `INSERT INTO partnerships (id, lead_id, brand_id, mitra_id, leader_id, outlet_id, package_id, progress_percentage, status, start_date, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.LeadID, p.BrandID, p.MitraID, p.LeaderID, p.OutletID, p.PackageID,
		p.ProgressPercentage, p.Status, p.StartDate, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PartnershipRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Partnership, error) {
	p := &entity.Partnership{}
	p.Mitra = &entity.User{}

	query := `SELECT p.id, p.lead_id, p.brand_id, p.mitra_id, p.leader_id, p.outlet_id, p.package_id,
			  p.progress_percentage, p.status, p.start_date, p.created_at, p.updated_at,
			  u.id, u.name, u.email, u.phone,
			  lu.id, lu.name, lu.email,
			  o.name,
			  op.name, op.price
			  FROM partnerships p
			  JOIN users u ON p.mitra_id = u.id
			  LEFT JOIN users lu ON p.leader_id = lu.id
			  LEFT JOIN outlets o ON p.outlet_id = o.id
			  LEFT JOIN outlet_packages op ON p.package_id = op.id
			  WHERE p.id = $1 AND p.deleted_at IS NULL`

	var leaderID sql.NullString
	var leaderName, leaderEmail sql.NullString
	var outletName, pkgName sql.NullString
	var pkgPrice sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.LeadID, &p.BrandID, &p.MitraID, &p.LeaderID, &p.OutletID, &p.PackageID,
		&p.ProgressPercentage, &p.Status, &p.StartDate, &p.CreatedAt, &p.UpdatedAt,
		&p.Mitra.ID, &p.Mitra.Name, &p.Mitra.Email, &p.Mitra.Phone,
		&leaderID, &leaderName, &leaderEmail,
		&outletName,
		&pkgName, &pkgPrice,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("partnership not found")
	}
	if err != nil {
		return nil, err
	}
	if leaderName.Valid {
		p.Leader = &entity.User{Name: leaderName.String, Email: leaderEmail.String}
		if leaderID.Valid {
			uid, _ := uuid.Parse(leaderID.String)
			p.Leader.ID = uid
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

	query := `SELECT p.id, p.lead_id, p.brand_id, p.mitra_id, p.leader_id, p.outlet_id, p.package_id,
			  p.progress_percentage, p.status, p.start_date, p.created_at, p.updated_at,
			  u.id, u.name, u.email,
			  lu.id, lu.name, lu.email,
			  o.name, op.name, op.price
			  FROM partnerships p
			  JOIN users u ON p.mitra_id = u.id
			  LEFT JOIN users lu ON p.leader_id = lu.id
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
		var leaderID sql.NullString
		var leaderName, leaderEmail sql.NullString
		var outletName, pkgName sql.NullString
		var pkgPrice sql.NullInt64
		err := rows.Scan(
			&p.ID, &p.LeadID, &p.BrandID, &p.MitraID, &p.LeaderID, &p.OutletID, &p.PackageID,
			&p.ProgressPercentage, &p.Status, &p.StartDate, &p.CreatedAt, &p.UpdatedAt,
			&p.Mitra.ID, &p.Mitra.Name, &p.Mitra.Email,
			&leaderID, &leaderName, &leaderEmail,
			&outletName, &pkgName, &pkgPrice,
		)
		if err != nil {
			return nil, 0, err
		}
		if leaderName.Valid {
			p.Leader = &entity.User{Name: leaderName.String, Email: leaderEmail.String}
			if leaderID.Valid {
				uid, _ := uuid.Parse(leaderID.String)
				p.Leader.ID = uid
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
	query := `SELECT p.id, p.lead_id, p.brand_id, p.mitra_id, p.progress_percentage, p.status, p.start_date, p.created_at, p.updated_at
			  FROM partnerships p
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
		err := rows.Scan(
			&p.ID, &p.LeadID, &p.BrandID, &p.MitraID, &p.ProgressPercentage,
			&p.Status, &p.StartDate, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
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
