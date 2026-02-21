package postgres

import (
	"context"
	"database/sql"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OutletPackageRepo struct {
	db *sql.DB
}

func NewOutletPackageRepo(db *sql.DB) *OutletPackageRepo {
	return &OutletPackageRepo{db: db}
}

const pkgColumns = `id, outlet_id, name, slug, price, minimum_dp, duration, image, estimated_bep, net_profit, description, benefits, sort_order, is_active, created_at, updated_at`

func (r *OutletPackageRepo) Create(ctx context.Context, pkg *entity.OutletPackage) error {
	query := `INSERT INTO outlet_packages (` + pkgColumns + `)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	_, err := r.db.ExecContext(ctx, query,
		pkg.ID, pkg.OutletID, pkg.Name, pkg.Slug, pkg.Price, pkg.MinimumDP, pkg.Duration,
		pkg.Image, pkg.EstimatedBEP, pkg.NetProfit, pkg.Description,
		pq.Array(pkg.Benefits), pkg.SortOrder, pkg.IsActive, pkg.CreatedAt, pkg.UpdatedAt,
	)
	return err
}

func scanPackage(scanner interface{ Scan(...any) error }) (*entity.OutletPackage, error) {
	pkg := &entity.OutletPackage{}
	var duration, image, estimatedBEP, netProfit, description sql.NullString
	err := scanner.Scan(
		&pkg.ID, &pkg.OutletID, &pkg.Name, &pkg.Slug, &pkg.Price, &pkg.MinimumDP,
		&duration, &image, &estimatedBEP, &netProfit, &description,
		&pkg.Benefits, &pkg.SortOrder, &pkg.IsActive, &pkg.CreatedAt, &pkg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	pkg.Duration = duration.String
	pkg.Image = image.String
	pkg.EstimatedBEP = estimatedBEP.String
	pkg.NetProfit = netProfit.String
	pkg.Description = description.String
	return pkg, nil
}

func (r *OutletPackageRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.OutletPackage, error) {
	query := `SELECT ` + pkgColumns + ` FROM outlet_packages WHERE id = $1`
	pkg, err := scanPackage(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return pkg, err
}

func (r *OutletPackageRepo) FindByOutletID(ctx context.Context, outletID uuid.UUID) ([]*entity.OutletPackage, error) {
	query := `SELECT ` + pkgColumns + ` FROM outlet_packages WHERE outlet_id = $1 ORDER BY sort_order ASC, created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, outletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []*entity.OutletPackage
	for rows.Next() {
		pkg, err := scanPackage(rows)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

func (r *OutletPackageRepo) Update(ctx context.Context, pkg *entity.OutletPackage) error {
	query := `UPDATE outlet_packages SET name=$1, slug=$2, price=$3, minimum_dp=$4, duration=$5, image=$6, estimated_bep=$7, net_profit=$8, description=$9, benefits=$10, sort_order=$11, is_active=$12, updated_at=$13
		WHERE id=$14`
	_, err := r.db.ExecContext(ctx, query,
		pkg.Name, pkg.Slug, pkg.Price, pkg.MinimumDP, pkg.Duration, pkg.Image,
		pkg.EstimatedBEP, pkg.NetProfit, pkg.Description,
		pq.Array(pkg.Benefits), pkg.SortOrder, pkg.IsActive, pkg.UpdatedAt, pkg.ID,
	)
	return err
}

func (r *OutletPackageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM outlet_packages WHERE id = $1`, id)
	return err
}
