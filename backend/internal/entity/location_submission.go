package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	LocSubStatusDraft     = "DRAFT"
	LocSubStatusSubmitted = "SUBMITTED"
	LocSubStatusInReview  = "IN_REVIEW"
	LocSubStatusScheduled = "SURVEY_SCHEDULED"
	LocSubStatusSurveyed  = "SURVEYED"
	LocSubStatusApproved  = "APPROVED"
	LocSubStatusRejected  = "REJECTED"
	LocSubStatusRevision  = "REVISION_NEEDED"
)

type LocationSubmission struct {
	ID            uuid.UUID  `json:"id"`
	MitraID       uuid.UUID  `json:"mitra_id"`
	PartnershipID *uuid.UUID `json:"partnership_id,omitempty"`
	NamaLokasi    string     `json:"nama_lokasi"`
	Alamat        string     `json:"alamat"`
	Provinsi      string     `json:"provinsi"`
	Kota          string     `json:"kota"`
	Kecamatan     string     `json:"kecamatan"`
	KodePos       string     `json:"kode_pos"`

	LuasTempat   float64 `json:"luas_tempat"`
	HargaSewa    int64   `json:"harga_sewa_per_tahun"`
	DurasiSewa   int     `json:"durasi_sewa"`
	TipeBangunan string  `json:"tipe_bangunan"`
	LebarJalan   float64 `json:"lebar_jalan"`
	JumlahLantai int     `json:"jumlah_lantai"`

	EstimasiLaluLintas int    `json:"estimasi_lalu_lintas"`
	DekatDengan        string `json:"dekat_dengan"`
	JumlahKompetitor   int    `json:"jumlah_kompetitor"`
	TargetMarket       string `json:"target_market"`

	ScoreTraffic    int    `json:"score_traffic"`
	ScoreSewa       int    `json:"score_sewa"`
	ScoreKompetitor int    `json:"score_kompetitor"`
	ScoreAkses      int    `json:"score_akses"`
	ScoreMarket     int    `json:"score_market"`
	TotalScore      int    `json:"total_score"`
	ScoreCategory   string `json:"score_category"`

	Status  string `json:"status"`
	Catatan string `json:"catatan"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Mitra       *User              `json:"mitra,omitempty"`
	Partnership *Partnership       `json:"partnership,omitempty"`
	Surveys     []LocationSurvey   `json:"surveys,omitempty"`
	Files       []LocationFile     `json:"files,omitempty"`
	Approvals   []LocationApproval `json:"approvals,omitempty"`
}

type LocationSurvey struct {
	ID            uuid.UUID  `json:"id"`
	LocationID    uuid.UUID  `json:"location_id"`
	SurveyDate    *time.Time `json:"survey_date,omitempty"`
	SurveyBy      *uuid.UUID `json:"survey_by,omitempty"`
	HasilSurvey   string     `json:"hasil_survey"`
	CatatanSurvey string     `json:"catatan_survey"`
	EstimasiOmzet int64      `json:"estimasi_omzet"`
	EstimasiBEP   int        `json:"estimasi_bep"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	Surveyor *User `json:"surveyor,omitempty"`
}

type LocationFile struct {
	ID         uuid.UUID  `json:"id"`
	LocationID uuid.UUID  `json:"location_id"`
	FileURL    string     `json:"file_url"`
	FileType   string     `json:"file_type"`
	Label      string     `json:"label"`
	UploadedBy *uuid.UUID `json:"uploaded_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type LocationApproval struct {
	ID         uuid.UUID `json:"id"`
	LocationID uuid.UUID `json:"location_id"`
	ApprovedBy uuid.UUID `json:"approved_by"`
	Decision   string    `json:"decision"`
	Note       string    `json:"note"`
	ApprovedAt time.Time `json:"approved_at"`

	Approver *User `json:"approver,omitempty"`
}
