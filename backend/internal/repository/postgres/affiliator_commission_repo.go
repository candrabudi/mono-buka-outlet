package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type AffiliatorCommissionRepo struct {
	db *sql.DB
}

func NewAffiliatorCommissionRepo(db *sql.DB) *AffiliatorCommissionRepo {
	return &AffiliatorCommissionRepo{db: db}
}

func (r *AffiliatorCommissionRepo) Create(ctx context.Context, c *entity.AffiliatorCommission) error {
	query := `
		INSERT INTO affiliator_commissions (id, affiliator_id, partnership_id, amount, type, description, given_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	var partnershipID interface{}
	if c.PartnershipID != nil {
		partnershipID = *c.PartnershipID
	}
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.AffiliatorID, partnershipID, c.Amount, c.Type, c.Description, c.GivenBy, c.CreatedAt,
	)
	return err
}

func (r *AffiliatorCommissionRepo) FindByAffiliatorID(ctx context.Context, affiliatorID uuid.UUID, page, limit int) ([]*entity.AffiliatorCommission, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM affiliator_commissions WHERE affiliator_id = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, affiliatorID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT c.id, c.affiliator_id, c.partnership_id, c.amount, c.type, c.description, c.given_by, c.created_at,
		       u.name,
		       COALESCE(mitra.name, ''),
		       COALESCE(o.name, '')
		FROM affiliator_commissions c
		JOIN users u ON u.id = c.given_by
		LEFT JOIN partnerships p ON p.id = c.partnership_id
		LEFT JOIN users mitra ON mitra.id = p.mitra_id
		LEFT JOIN outlets o ON o.id = p.outlet_id
		WHERE c.affiliator_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	offset := (page - 1) * limit
	rows, err := r.db.QueryContext(ctx, query, affiliatorID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*entity.AffiliatorCommission
	for rows.Next() {
		c := &entity.AffiliatorCommission{}
		c.GivenByUser = &entity.User{}
		var partnershipID sql.NullString
		var description sql.NullString
		var mitraName string
		var outletName string

		err := rows.Scan(
			&c.ID, &c.AffiliatorID, &partnershipID, &c.Amount, &c.Type, &description, &c.GivenBy, &c.CreatedAt,
			&c.GivenByUser.Name,
			&mitraName,
			&outletName,
		)
		if err != nil {
			return nil, 0, err
		}
		if partnershipID.Valid {
			pid, _ := uuid.Parse(partnershipID.String)
			c.PartnershipID = &pid
			c.Partnership = &entity.Partnership{
				ID: pid,
				Mitra:  &entity.User{Name: mitraName},
				Outlet: &entity.Outlet{Name: outletName},
			}
		}
		if description.Valid {
			c.Description = description.String
		}
		results = append(results, c)
	}
	return results, total, nil
}

func (r *AffiliatorCommissionRepo) GetBalance(ctx context.Context, affiliatorID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(
			(SELECT SUM(amount) FROM affiliator_commissions WHERE affiliator_id = $1), 0
		) - COALESCE(
			(SELECT SUM(amount) FROM affiliator_withdrawals WHERE affiliator_id = $1 AND status IN ('APPROVED', 'TRANSFERRED')), 0
		)
	`
	var balance float64
	err := r.db.QueryRowContext(ctx, query, affiliatorID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil saldo: %w", err)
	}
	return balance, nil
}

func (r *AffiliatorCommissionRepo) GetTotalEarned(ctx context.Context, affiliatorID uuid.UUID) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM affiliator_commissions WHERE affiliator_id = $1`
	var total float64
	err := r.db.QueryRowContext(ctx, query, affiliatorID).Scan(&total)
	return total, err
}
