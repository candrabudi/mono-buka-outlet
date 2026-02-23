package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

// ══════════════════════════════════════
// EbookRepo
// ══════════════════════════════════════

type EbookRepo struct {
	db *sql.DB
}

func NewEbookRepo(db *sql.DB) *EbookRepo {
	return &EbookRepo{db: db}
}

func (r *EbookRepo) Create(ctx context.Context, ebook *entity.Ebook) error {
	if ebook.ID == uuid.Nil {
		ebook.ID = uuid.New()
	}
	ebook.CreatedAt = time.Now()
	ebook.UpdatedAt = time.Now()

	query := `INSERT INTO ebooks (id, title, slug, description, author, cover_url, file_url, price, is_active, created_at, updated_at)
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.db.ExecContext(ctx, query,
		ebook.ID, ebook.Title, ebook.Slug, ebook.Description, ebook.Author,
		ebook.CoverURL, ebook.FileURL, ebook.Price, ebook.IsActive,
		ebook.CreatedAt, ebook.UpdatedAt)
	return err
}

func (r *EbookRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Ebook, error) {
	ebook := &entity.Ebook{}
	query := `SELECT id, title, slug, description, author, cover_url, file_url, price, is_active, total_sold, created_at, updated_at
			  FROM ebooks WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ebook.ID, &ebook.Title, &ebook.Slug, &ebook.Description, &ebook.Author,
		&ebook.CoverURL, &ebook.FileURL, &ebook.Price, &ebook.IsActive, &ebook.TotalSold,
		&ebook.CreatedAt, &ebook.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ebook not found")
	}
	if err != nil {
		return nil, err
	}
	ebook.Categories, _ = r.loadCategories(ctx, ebook.ID)
	return ebook, nil
}

func (r *EbookRepo) FindBySlug(ctx context.Context, slug string) (*entity.Ebook, error) {
	ebook := &entity.Ebook{}
	query := `SELECT id, title, slug, description, author, cover_url, file_url, price, is_active, total_sold, created_at, updated_at
			  FROM ebooks WHERE slug = $1`
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&ebook.ID, &ebook.Title, &ebook.Slug, &ebook.Description, &ebook.Author,
		&ebook.CoverURL, &ebook.FileURL, &ebook.Price, &ebook.IsActive, &ebook.TotalSold,
		&ebook.CreatedAt, &ebook.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ebook not found")
	}
	if err != nil {
		return nil, err
	}
	ebook.Categories, _ = r.loadCategories(ctx, ebook.ID)
	return ebook, nil
}

func (r *EbookRepo) FindAll(ctx context.Context, activeOnly bool, search string, page, limit int) ([]*entity.Ebook, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if activeOnly {
		where = append(where, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, true)
		argIdx++
	}
	if search != "" {
		where = append(where, fmt.Sprintf(`(LOWER(title) LIKE $%d OR LOWER(author) LIKE $%d
			OR EXISTS (SELECT 1 FROM ebook_category_mapping ecm
				JOIN ebook_categories ec ON ec.id = ecm.category_id
				WHERE ecm.ebook_id = ebooks.id AND LOWER(ec.name) LIKE $%d))`, argIdx, argIdx, argIdx))
		args = append(args, "%"+strings.ToLower(search)+"%")
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ebooks WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	dataQuery := fmt.Sprintf(`SELECT id, title, slug, description, author, cover_url, file_url, price, is_active, total_sold, created_at, updated_at
		FROM ebooks WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var ebooks []*entity.Ebook
	for rows.Next() {
		e := &entity.Ebook{}
		if err := rows.Scan(&e.ID, &e.Title, &e.Slug, &e.Description, &e.Author,
			&e.CoverURL, &e.FileURL, &e.Price, &e.IsActive, &e.TotalSold,
			&e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, 0, err
		}
		e.Categories, _ = r.loadCategories(ctx, e.ID)
		ebooks = append(ebooks, e)
	}
	return ebooks, total, nil
}

func (r *EbookRepo) Update(ctx context.Context, ebook *entity.Ebook) error {
	ebook.UpdatedAt = time.Now()
	query := `UPDATE ebooks SET title=$1, slug=$2, description=$3, author=$4, cover_url=$5, file_url=$6, price=$7, is_active=$8, updated_at=$9 WHERE id=$10`
	_, err := r.db.ExecContext(ctx, query,
		ebook.Title, ebook.Slug, ebook.Description, ebook.Author,
		ebook.CoverURL, ebook.FileURL, ebook.Price, ebook.IsActive,
		ebook.UpdatedAt, ebook.ID)
	return err
}

func (r *EbookRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM ebooks WHERE id = $1", id)
	return err
}

func (r *EbookRepo) IncrementSold(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE ebooks SET total_sold = total_sold + 1, updated_at = NOW() WHERE id = $1", id)
	return err
}

// SyncCategories replaces all category mappings for an ebook
func (r *EbookRepo) SyncCategories(ctx context.Context, ebookID uuid.UUID, categoryIDs []uuid.UUID) error {
	// Delete old mappings
	_, err := r.db.ExecContext(ctx, "DELETE FROM ebook_category_mapping WHERE ebook_id = $1", ebookID)
	if err != nil {
		return err
	}
	// Insert new ones
	for _, catID := range categoryIDs {
		_, err := r.db.ExecContext(ctx,
			"INSERT INTO ebook_category_mapping (ebook_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			ebookID, catID)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadCategories loads categories for an ebook via the pivot table
func (r *EbookRepo) loadCategories(ctx context.Context, ebookID uuid.UUID) ([]*entity.EbookCategory, error) {
	query := `SELECT ec.id, ec.name, ec.slug, ec.is_active, ec.created_at, ec.updated_at
		FROM ebook_categories ec
		JOIN ebook_category_mapping ecm ON ecm.category_id = ec.id
		WHERE ecm.ebook_id = $1 ORDER BY ec.name`
	rows, err := r.db.QueryContext(ctx, query, ebookID)
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

// ══════════════════════════════════════
// EbookOrderRepo
// ══════════════════════════════════════

type EbookOrderRepo struct {
	db *sql.DB
}

func NewEbookOrderRepo(db *sql.DB) *EbookOrderRepo {
	return &EbookOrderRepo{db: db}
}

var ebookOrderColumns = `eo.id, eo.ebook_id, eo.user_id, eo.order_number, eo.amount, eo.status,
	eo.download_status, eo.download_requested_at, eo.download_approved_at, COALESCE(eo.download_note,''),
	COALESCE(eo.midtrans_order_id,''), COALESCE(eo.midtrans_snap_token,''), COALESCE(eo.midtrans_redirect_url,''),
	COALESCE(eo.midtrans_payment_type,''), COALESCE(eo.midtrans_transaction_id,''), COALESCE(eo.midtrans_transaction_status,''),
	eo.paid_at, eo.expired_at, COALESCE(eo.payment_proof_url,''), eo.created_at, eo.updated_at`

func scanEbookOrder(row interface {
	Scan(dest ...interface{}) error
}) (*entity.EbookOrder, error) {
	o := &entity.EbookOrder{}
	err := row.Scan(
		&o.ID, &o.EbookID, &o.UserID, &o.OrderNumber, &o.Amount, &o.Status,
		&o.DownloadStatus, &o.DownloadRequestedAt, &o.DownloadApprovedAt, &o.DownloadNote,
		&o.MidtransOrderID, &o.MidtransSnapToken, &o.MidtransRedirectURL,
		&o.MidtransPaymentType, &o.MidtransTransactionID, &o.MidtransTransactionStatus,
		&o.PaidAt, &o.ExpiredAt, &o.PaymentProofURL, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *EbookOrderRepo) Create(ctx context.Context, order *entity.EbookOrder) error {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	if order.DownloadStatus == "" {
		order.DownloadStatus = entity.DownloadStatusNone
	}

	query := `INSERT INTO ebook_orders (id, ebook_id, user_id, order_number, amount, status, download_status,
			  midtrans_order_id, midtrans_snap_token, midtrans_redirect_url, expired_at, created_at, updated_at)
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.EbookID, order.UserID, order.OrderNumber, order.Amount, order.Status, order.DownloadStatus,
		order.MidtransOrderID, order.MidtransSnapToken, order.MidtransRedirectURL,
		order.ExpiredAt, order.CreatedAt, order.UpdatedAt)
	return err
}

func (r *EbookOrderRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.EbookOrder, error) {
	query := fmt.Sprintf(`SELECT %s, e.title, e.slug, e.cover_url, e.price, u.name, u.email
		FROM ebook_orders eo
		LEFT JOIN ebooks e ON e.id = eo.ebook_id
		LEFT JOIN users u ON u.id = eo.user_id
		WHERE eo.id = $1`, ebookOrderColumns)
	row := r.db.QueryRowContext(ctx, query, id)

	o := &entity.EbookOrder{}
	var ebookTitle, ebookSlug, ebookCover string
	var ebookPrice int64
	var userName, userEmail string
	err := row.Scan(
		&o.ID, &o.EbookID, &o.UserID, &o.OrderNumber, &o.Amount, &o.Status,
		&o.DownloadStatus, &o.DownloadRequestedAt, &o.DownloadApprovedAt, &o.DownloadNote,
		&o.MidtransOrderID, &o.MidtransSnapToken, &o.MidtransRedirectURL,
		&o.MidtransPaymentType, &o.MidtransTransactionID, &o.MidtransTransactionStatus,
		&o.PaidAt, &o.ExpiredAt, &o.PaymentProofURL, &o.CreatedAt, &o.UpdatedAt,
		&ebookTitle, &ebookSlug, &ebookCover, &ebookPrice, &userName, &userEmail)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ebook order not found")
	}
	if err != nil {
		return nil, err
	}
	o.Ebook = &entity.Ebook{ID: o.EbookID, Title: ebookTitle, Slug: ebookSlug, CoverURL: ebookCover, Price: ebookPrice}
	o.User = &entity.User{ID: o.UserID, Name: userName, Email: userEmail}
	return o, nil
}

func (r *EbookOrderRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.EbookOrder, error) {
	query := fmt.Sprintf(`SELECT %s, e.title, e.slug, e.cover_url, e.price
		FROM ebook_orders eo
		LEFT JOIN ebooks e ON e.id = eo.ebook_id
		WHERE eo.user_id = $1 ORDER BY eo.created_at DESC`, ebookOrderColumns)
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.EbookOrder
	for rows.Next() {
		o := &entity.EbookOrder{}
		var ebookTitle, ebookSlug, ebookCover string
		var ebookPrice int64
		if err := rows.Scan(
			&o.ID, &o.EbookID, &o.UserID, &o.OrderNumber, &o.Amount, &o.Status,
			&o.DownloadStatus, &o.DownloadRequestedAt, &o.DownloadApprovedAt, &o.DownloadNote,
			&o.MidtransOrderID, &o.MidtransSnapToken, &o.MidtransRedirectURL,
			&o.MidtransPaymentType, &o.MidtransTransactionID, &o.MidtransTransactionStatus,
			&o.PaidAt, &o.ExpiredAt, &o.PaymentProofURL, &o.CreatedAt, &o.UpdatedAt,
			&ebookTitle, &ebookSlug, &ebookCover, &ebookPrice); err != nil {
			return nil, err
		}
		o.Ebook = &entity.Ebook{ID: o.EbookID, Title: ebookTitle, Slug: ebookSlug, CoverURL: ebookCover, Price: ebookPrice}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *EbookOrderRepo) FindByMidtransOrderID(ctx context.Context, orderID string) (*entity.EbookOrder, error) {
	query := fmt.Sprintf(`SELECT %s FROM ebook_orders eo WHERE eo.midtrans_order_id = $1`, ebookOrderColumns)
	o, err := scanEbookOrder(r.db.QueryRowContext(ctx, query, orderID))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ebook order not found for midtrans order: %s", orderID)
	}
	return o, err
}

func (r *EbookOrderRepo) HasUserPurchased(ctx context.Context, userID, ebookID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM ebook_orders WHERE user_id = $1 AND ebook_id = $2 AND status = $3`,
		userID, ebookID, entity.EbookOrderStatusPaid).Scan(&count)
	return count > 0, err
}

func (r *EbookOrderRepo) UpdateMidtransStatus(ctx context.Context, orderID string, txnStatus, paymentType, txnID string) error {
	now := time.Now()
	status := entity.EbookOrderStatusPending

	switch txnStatus {
	case "capture", "settlement":
		status = entity.EbookOrderStatusPaid
	case "expire":
		status = entity.EbookOrderStatusExpired
	case "cancel", "deny":
		status = entity.EbookOrderStatusFailed
	case "pending":
		status = entity.EbookOrderStatusPending
	}

	query := `UPDATE ebook_orders SET
			  midtrans_transaction_status = $1,
			  midtrans_payment_type = $2,
			  midtrans_transaction_id = $3,
			  status = $4,
			  updated_at = $5`
	args := []interface{}{txnStatus, paymentType, txnID, status, now}

	if status == entity.EbookOrderStatusPaid {
		query += `, paid_at = $6 WHERE midtrans_order_id = $7`
		args = append(args, now, orderID)
	} else {
		query += ` WHERE midtrans_order_id = $6`
		args = append(args, orderID)
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *EbookOrderRepo) GenerateOrderNumber(ctx context.Context) (string, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ebook_orders`).Scan(&count)
	if err != nil {
		return "", err
	}
	now := time.Now()
	return fmt.Sprintf("EB-%s-%04d", now.Format("20060102"), count+1), nil
}

func (r *EbookOrderRepo) RequestDownload(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE ebook_orders SET download_status = $1, download_requested_at = $2, updated_at = $3
			  WHERE id = $4 AND status = $5 AND download_status = $6`
	res, err := r.db.ExecContext(ctx, query, entity.DownloadStatusRequested, now, now, id, entity.EbookOrderStatusPaid, entity.DownloadStatusNone)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order not eligible for download request")
	}
	return nil
}

func (r *EbookOrderRepo) ApproveDownload(ctx context.Context, id uuid.UUID, note string) error {
	now := time.Now()
	query := `UPDATE ebook_orders SET download_status = $1, download_approved_at = $2, download_note = $3, updated_at = $4
			  WHERE id = $5 AND download_status = $6`
	res, err := r.db.ExecContext(ctx, query, entity.DownloadStatusApproved, now, note, now, id, entity.DownloadStatusRequested)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order not found or not in REQUESTED status")
	}
	return nil
}

func (r *EbookOrderRepo) RejectDownload(ctx context.Context, id uuid.UUID, note string) error {
	now := time.Now()
	query := `UPDATE ebook_orders SET download_status = $1, download_note = $2, updated_at = $3
			  WHERE id = $4 AND download_status = $5`
	res, err := r.db.ExecContext(ctx, query, entity.DownloadStatusRejected, note, now, id, entity.DownloadStatusRequested)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order not found or not in REQUESTED status")
	}
	return nil
}

func (r *EbookOrderRepo) FindPendingDownloads(ctx context.Context) ([]*entity.EbookOrder, error) {
	query := fmt.Sprintf(`SELECT %s, e.title, e.slug, e.cover_url, e.price, u.name, u.email
		FROM ebook_orders eo
		LEFT JOIN ebooks e ON e.id = eo.ebook_id
		LEFT JOIN users u ON u.id = eo.user_id
		WHERE eo.download_status = $1
		ORDER BY eo.download_requested_at ASC`, ebookOrderColumns)
	rows, err := r.db.QueryContext(ctx, query, entity.DownloadStatusRequested)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.EbookOrder
	for rows.Next() {
		o := &entity.EbookOrder{}
		var ebookTitle, ebookSlug, ebookCover string
		var ebookPrice int64
		var userName, userEmail string
		if err := rows.Scan(
			&o.ID, &o.EbookID, &o.UserID, &o.OrderNumber, &o.Amount, &o.Status,
			&o.DownloadStatus, &o.DownloadRequestedAt, &o.DownloadApprovedAt, &o.DownloadNote,
			&o.MidtransOrderID, &o.MidtransSnapToken, &o.MidtransRedirectURL,
			&o.MidtransPaymentType, &o.MidtransTransactionID, &o.MidtransTransactionStatus,
			&o.PaidAt, &o.ExpiredAt, &o.CreatedAt, &o.UpdatedAt,
			&ebookTitle, &ebookSlug, &ebookCover, &ebookPrice, &userName, &userEmail); err != nil {
			return nil, err
		}
		o.Ebook = &entity.Ebook{ID: o.EbookID, Title: ebookTitle, Slug: ebookSlug, CoverURL: ebookCover, Price: ebookPrice}
		o.User = &entity.User{ID: o.UserID, Name: userName, Email: userEmail}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *EbookOrderRepo) FindAllOrders(ctx context.Context, status, downloadStatus string, page, limit int) ([]*entity.EbookOrder, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		where = append(where, fmt.Sprintf("eo.status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}
	if downloadStatus != "" {
		where = append(where, fmt.Sprintf("eo.download_status = $%d", argIdx))
		args = append(args, downloadStatus)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM ebook_orders eo WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	dataQ := fmt.Sprintf(`SELECT %s, e.title, e.slug, e.cover_url, e.price, u.name, u.email
		FROM ebook_orders eo
		LEFT JOIN ebooks e ON e.id = eo.ebook_id
		LEFT JOIN users u ON u.id = eo.user_id
		WHERE %s ORDER BY eo.created_at DESC LIMIT $%d OFFSET $%d`,
		ebookOrderColumns, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQ, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*entity.EbookOrder
	for rows.Next() {
		o := &entity.EbookOrder{}
		var ebookTitle, ebookSlug, ebookCover string
		var ebookPrice int64
		var userName, userEmail string
		if err := rows.Scan(
			&o.ID, &o.EbookID, &o.UserID, &o.OrderNumber, &o.Amount, &o.Status,
			&o.DownloadStatus, &o.DownloadRequestedAt, &o.DownloadApprovedAt, &o.DownloadNote,
			&o.MidtransOrderID, &o.MidtransSnapToken, &o.MidtransRedirectURL,
			&o.MidtransPaymentType, &o.MidtransTransactionID, &o.MidtransTransactionStatus,
			&o.PaidAt, &o.ExpiredAt, &o.PaymentProofURL, &o.CreatedAt, &o.UpdatedAt,
			&ebookTitle, &ebookSlug, &ebookCover, &ebookPrice, &userName, &userEmail); err != nil {
			return nil, 0, err
		}
		o.Ebook = &entity.Ebook{ID: o.EbookID, Title: ebookTitle, Slug: ebookSlug, CoverURL: ebookCover, Price: ebookPrice}
		o.User = &entity.User{ID: o.UserID, Name: userName, Email: userEmail}
		orders = append(orders, o)
	}
	return orders, total, nil
}

func (r *EbookOrderRepo) CancelOrder(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE ebook_orders SET status = $1, updated_at = NOW() WHERE id = $2 AND status = $3`,
		entity.EbookOrderStatusCanceled, id, entity.EbookOrderStatusPending)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("pesanan tidak ditemukan atau tidak bisa dibatalkan")
	}
	return nil
}

func (r *EbookOrderRepo) UploadPaymentProof(ctx context.Context, id uuid.UUID, proofURL string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE ebook_orders SET payment_proof_url = $1, updated_at = NOW() WHERE id = $2 AND status = $3`,
		proofURL, id, entity.EbookOrderStatusPending)
	return err
}

func (r *EbookOrderRepo) ApprovePayment(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE ebook_orders SET status = $1, paid_at = $2, midtrans_payment_type = 'manual_transfer', updated_at = NOW() WHERE id = $3 AND status = $4`,
		entity.EbookOrderStatusPaid, now, id, entity.EbookOrderStatusPending)
	return err
}

func (r *EbookOrderRepo) RejectPayment(ctx context.Context, id uuid.UUID, note string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE ebook_orders SET status = $1, download_note = $2, updated_at = NOW() WHERE id = $3 AND status = $4`,
		entity.EbookOrderStatusFailed, note, id, entity.EbookOrderStatusPending)
	return err
}
