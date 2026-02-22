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

type LocationSubmissionRepo struct {
	db *sql.DB
}

func NewLocationSubmissionRepo(db *sql.DB) *LocationSubmissionRepo {
	return &LocationSubmissionRepo{db: db}
}

// ---------- SUBMISSION CRUD ----------

func (r *LocationSubmissionRepo) Create(ctx context.Context, ls *entity.LocationSubmission) error {
	if ls.ID == uuid.Nil {
		ls.ID = uuid.New()
	}
	ls.CreatedAt = time.Now()
	ls.UpdatedAt = time.Now()
	if ls.Status == "" {
		ls.Status = entity.LocSubStatusDraft
	}

	query := `INSERT INTO location_submissions (
		id, mitra_id, partnership_id, nama_lokasi, alamat, provinsi, kota, kecamatan, kode_pos,
		luas_tempat, harga_sewa_per_tahun, durasi_sewa,
		tipe_bangunan, lebar_jalan, jumlah_lantai,
		estimasi_lalu_lintas, dekat_dengan, jumlah_kompetitor, target_market,
		score_traffic, score_sewa, score_kompetitor, score_akses, score_market,
		total_score, score_category, status, catatan, created_at, updated_at
	) VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
		$21,$22,$23,$24,$25,$26,$27,$28,$29,$30
	)`
	_, err := r.db.ExecContext(ctx, query,
		ls.ID, ls.MitraID, ls.PartnershipID, ls.NamaLokasi, ls.Alamat, ls.Provinsi, ls.Kota, ls.Kecamatan, ls.KodePos,
		ls.LuasTempat, ls.HargaSewa, ls.DurasiSewa,
		ls.TipeBangunan, ls.LebarJalan, ls.JumlahLantai,
		ls.EstimasiLaluLintas, ls.DekatDengan, ls.JumlahKompetitor, ls.TargetMarket,
		ls.ScoreTraffic, ls.ScoreSewa, ls.ScoreKompetitor, ls.ScoreAkses, ls.ScoreMarket,
		ls.TotalScore, ls.ScoreCategory, ls.Status, ls.Catatan, ls.CreatedAt, ls.UpdatedAt,
	)
	return err
}

func (r *LocationSubmissionRepo) Update(ctx context.Context, ls *entity.LocationSubmission) error {
	ls.UpdatedAt = time.Now()
	query := `UPDATE location_submissions SET
		partnership_id=$1, nama_lokasi=$2, alamat=$3, provinsi=$4, kota=$5, kecamatan=$6, kode_pos=$7,
		luas_tempat=$8, harga_sewa_per_tahun=$9, durasi_sewa=$10,
		tipe_bangunan=$11, lebar_jalan=$12, jumlah_lantai=$13,
		estimasi_lalu_lintas=$14, dekat_dengan=$15, jumlah_kompetitor=$16, target_market=$17,
		score_traffic=$18, score_sewa=$19, score_kompetitor=$20, score_akses=$21, score_market=$22,
		total_score=$23, score_category=$24, status=$25, catatan=$26, updated_at=$27
		WHERE id=$28`
	_, err := r.db.ExecContext(ctx, query,
		ls.PartnershipID, ls.NamaLokasi, ls.Alamat, ls.Provinsi, ls.Kota, ls.Kecamatan, ls.KodePos,
		ls.LuasTempat, ls.HargaSewa, ls.DurasiSewa,
		ls.TipeBangunan, ls.LebarJalan, ls.JumlahLantai,
		ls.EstimasiLaluLintas, ls.DekatDengan, ls.JumlahKompetitor, ls.TargetMarket,
		ls.ScoreTraffic, ls.ScoreSewa, ls.ScoreKompetitor, ls.ScoreAkses, ls.ScoreMarket,
		ls.TotalScore, ls.ScoreCategory, ls.Status, ls.Catatan, ls.UpdatedAt, ls.ID,
	)
	return err
}

func (r *LocationSubmissionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE location_submissions SET status=$1, updated_at=$2 WHERE id=$3`,
		status, time.Now(), id)
	return err
}

func (r *LocationSubmissionRepo) UpdateScore(ctx context.Context, id uuid.UUID, traffic, sewa, kompetitor, akses, market, total int, category string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE location_submissions SET score_traffic=$1, score_sewa=$2, score_kompetitor=$3, score_akses=$4, score_market=$5, total_score=$6, score_category=$7, updated_at=$8 WHERE id=$9`,
		traffic, sewa, kompetitor, akses, market, total, category, time.Now(), id)
	return err
}

const locSubmissionCols = `ls.id, ls.mitra_id, ls.partnership_id, ls.nama_lokasi, ls.alamat, ls.provinsi, ls.kota, ls.kecamatan, ls.kode_pos,
	ls.luas_tempat, ls.harga_sewa_per_tahun, ls.durasi_sewa,
	ls.tipe_bangunan, ls.lebar_jalan, ls.jumlah_lantai,
	ls.estimasi_lalu_lintas, ls.dekat_dengan, ls.jumlah_kompetitor, ls.target_market,
	ls.score_traffic, ls.score_sewa, ls.score_kompetitor, ls.score_akses, ls.score_market,
	ls.total_score, COALESCE(ls.score_category,''), ls.status, COALESCE(ls.catatan,''),
	ls.created_at, ls.updated_at`

func scanSubmission(row interface{ Scan(...interface{}) error }) (*entity.LocationSubmission, error) {
	ls := &entity.LocationSubmission{}
	err := row.Scan(
		&ls.ID, &ls.MitraID, &ls.PartnershipID, &ls.NamaLokasi, &ls.Alamat, &ls.Provinsi, &ls.Kota, &ls.Kecamatan, &ls.KodePos,
		&ls.LuasTempat, &ls.HargaSewa, &ls.DurasiSewa,
		&ls.TipeBangunan, &ls.LebarJalan, &ls.JumlahLantai,
		&ls.EstimasiLaluLintas, &ls.DekatDengan, &ls.JumlahKompetitor, &ls.TargetMarket,
		&ls.ScoreTraffic, &ls.ScoreSewa, &ls.ScoreKompetitor, &ls.ScoreAkses, &ls.ScoreMarket,
		&ls.TotalScore, &ls.ScoreCategory, &ls.Status, &ls.Catatan,
		&ls.CreatedAt, &ls.UpdatedAt,
	)
	return ls, err
}

func (r *LocationSubmissionRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.LocationSubmission, error) {
	query := fmt.Sprintf(`SELECT %s FROM location_submissions ls WHERE ls.id=$1 AND ls.deleted_at IS NULL`, locSubmissionCols)
	ls, err := scanSubmission(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("location submission not found")
	}
	if err != nil {
		return nil, err
	}

	// Load mitra
	mitra := &entity.User{}
	_ = r.db.QueryRowContext(ctx, `SELECT id, name, email, phone FROM users WHERE id=$1`, ls.MitraID).Scan(&mitra.ID, &mitra.Name, &mitra.Email, &mitra.Phone)
	ls.Mitra = mitra

	// Load partnership
	if ls.PartnershipID != nil {
		p := &entity.Partnership{}
		_ = r.db.QueryRowContext(ctx, `SELECT id, mitra_id, status, created_at FROM partnerships WHERE id=$1`, *ls.PartnershipID).Scan(&p.ID, &p.MitraID, &p.Status, &p.CreatedAt)
		ls.Partnership = p
	}

	// Load surveys
	ls.Surveys, _ = r.FindSurveys(ctx, id)
	// Load files
	ls.Files, _ = r.FindFiles(ctx, id)
	// Load approvals
	ls.Approvals, _ = r.FindApprovals(ctx, id)

	return ls, nil
}

func (r *LocationSubmissionRepo) FindAll(ctx context.Context, status, kota, search string) ([]*entity.LocationSubmission, error) {
	where := []string{"ls.deleted_at IS NULL"}
	args := []interface{}{}
	n := 0

	if status != "" {
		n++
		where = append(where, fmt.Sprintf("ls.status=$%d", n))
		args = append(args, status)
	}
	if kota != "" {
		n++
		where = append(where, fmt.Sprintf("ls.kota ILIKE $%d", n))
		args = append(args, "%"+kota+"%")
	}
	if search != "" {
		n++
		where = append(where, fmt.Sprintf("(ls.nama_lokasi ILIKE $%d OR ls.alamat ILIKE $%d)", n, n))
		args = append(args, "%"+search+"%")
	}

	query := fmt.Sprintf(`SELECT %s FROM location_submissions ls WHERE %s ORDER BY ls.created_at DESC`,
		locSubmissionCols, strings.Join(where, " AND "))
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.LocationSubmission
	for rows.Next() {
		ls, err := scanSubmission(rows)
		if err != nil {
			return nil, err
		}
		// Load mitra name
		mitra := &entity.User{}
		_ = r.db.QueryRowContext(ctx, `SELECT id, name, email FROM users WHERE id=$1`, ls.MitraID).Scan(&mitra.ID, &mitra.Name, &mitra.Email)
		ls.Mitra = mitra
		results = append(results, ls)
	}
	return results, nil
}

func (r *LocationSubmissionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE location_submissions SET deleted_at=$1 WHERE id=$2`, time.Now(), id)
	return err
}

func (r *LocationSubmissionRepo) FindByPartnership(ctx context.Context, partnershipID uuid.UUID) ([]*entity.LocationSubmission, error) {
	query := fmt.Sprintf(`SELECT %s FROM location_submissions ls WHERE ls.partnership_id=$1 AND ls.deleted_at IS NULL ORDER BY ls.created_at DESC`, locSubmissionCols)
	rows, err := r.db.QueryContext(ctx, query, partnershipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.LocationSubmission
	for rows.Next() {
		ls, err := scanSubmission(rows)
		if err != nil {
			return nil, err
		}
		mitra := &entity.User{}
		_ = r.db.QueryRowContext(ctx, `SELECT id, name, email FROM users WHERE id=$1`, ls.MitraID).Scan(&mitra.ID, &mitra.Name, &mitra.Email)
		ls.Mitra = mitra
		results = append(results, ls)
	}
	return results, nil
}

func (r *LocationSubmissionRepo) FindByMitraID(ctx context.Context, mitraID uuid.UUID) ([]*entity.LocationSubmission, error) {
	query := fmt.Sprintf(`SELECT %s FROM location_submissions ls WHERE ls.mitra_id=$1 AND ls.deleted_at IS NULL ORDER BY ls.created_at DESC`, locSubmissionCols)
	rows, err := r.db.QueryContext(ctx, query, mitraID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.LocationSubmission
	for rows.Next() {
		ls, err := scanSubmission(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, ls)
	}
	if results == nil {
		results = []*entity.LocationSubmission{}
	}
	return results, nil
}

// ---------- SURVEYS ----------

func (r *LocationSubmissionRepo) CreateSurvey(ctx context.Context, s *entity.LocationSurvey) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	query := `INSERT INTO location_surveys (id, location_id, survey_date, survey_by, hasil_survey, catatan_survey, estimasi_omzet, estimasi_bep, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.LocationID, s.SurveyDate, s.SurveyBy, s.HasilSurvey, s.CatatanSurvey, s.EstimasiOmzet, s.EstimasiBEP, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *LocationSubmissionRepo) FindSurveys(ctx context.Context, locationID uuid.UUID) ([]entity.LocationSurvey, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT s.id, s.location_id, s.survey_date, s.survey_by, COALESCE(s.hasil_survey,''), COALESCE(s.catatan_survey,''),
		 s.estimasi_omzet, s.estimasi_bep, s.created_at, s.updated_at,
		 COALESCE(u.name,''), COALESCE(u.email,'')
		 FROM location_surveys s LEFT JOIN users u ON u.id=s.survey_by
		 WHERE s.location_id=$1 ORDER BY s.created_at DESC`, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []entity.LocationSurvey
	for rows.Next() {
		var s entity.LocationSurvey
		var uName, uEmail string
		if err := rows.Scan(&s.ID, &s.LocationID, &s.SurveyDate, &s.SurveyBy,
			&s.HasilSurvey, &s.CatatanSurvey, &s.EstimasiOmzet, &s.EstimasiBEP,
			&s.CreatedAt, &s.UpdatedAt, &uName, &uEmail); err != nil {
			return nil, err
		}
		if s.SurveyBy != nil {
			s.Surveyor = &entity.User{ID: *s.SurveyBy, Name: uName, Email: uEmail}
		}
		surveys = append(surveys, s)
	}
	return surveys, nil
}

// ---------- FILES ----------

func (r *LocationSubmissionRepo) CreateFile(ctx context.Context, f *entity.LocationFile) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	f.CreatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO location_files (id, location_id, file_url, file_type, label, uploaded_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		f.ID, f.LocationID, f.FileURL, f.FileType, f.Label, f.UploadedBy, f.CreatedAt)
	return err
}

func (r *LocationSubmissionRepo) FindFiles(ctx context.Context, locationID uuid.UUID) ([]entity.LocationFile, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, location_id, file_url, file_type, COALESCE(label,''), uploaded_by, created_at FROM location_files WHERE location_id=$1 ORDER BY created_at`, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []entity.LocationFile
	for rows.Next() {
		var f entity.LocationFile
		if err := rows.Scan(&f.ID, &f.LocationID, &f.FileURL, &f.FileType, &f.Label, &f.UploadedBy, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (r *LocationSubmissionRepo) DeleteFile(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM location_files WHERE id=$1`, id)
	return err
}

// ---------- APPROVALS ----------

func (r *LocationSubmissionRepo) CreateApproval(ctx context.Context, a *entity.LocationApproval) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	a.ApprovedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO location_approvals (id, location_id, approved_by, decision, note, approved_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		a.ID, a.LocationID, a.ApprovedBy, a.Decision, a.Note, a.ApprovedAt)
	return err
}

func (r *LocationSubmissionRepo) FindApprovals(ctx context.Context, locationID uuid.UUID) ([]entity.LocationApproval, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT a.id, a.location_id, a.approved_by, a.decision, COALESCE(a.note,''), a.approved_at,
		 COALESCE(u.name,''), COALESCE(u.email,'')
		 FROM location_approvals a LEFT JOIN users u ON u.id=a.approved_by
		 WHERE a.location_id=$1 ORDER BY a.approved_at DESC`, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvals []entity.LocationApproval
	for rows.Next() {
		var a entity.LocationApproval
		var uName, uEmail string
		if err := rows.Scan(&a.ID, &a.LocationID, &a.ApprovedBy, &a.Decision, &a.Note, &a.ApprovedAt, &uName, &uEmail); err != nil {
			return nil, err
		}
		a.Approver = &entity.User{ID: a.ApprovedBy, Name: uName, Email: uEmail}
		approvals = append(approvals, a)
	}
	return approvals, nil
}
