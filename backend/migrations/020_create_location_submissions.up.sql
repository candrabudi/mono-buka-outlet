-- Location Submissions
CREATE TABLE IF NOT EXISTS location_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mitra_id UUID NOT NULL REFERENCES users(id),
    nama_lokasi VARCHAR(255) NOT NULL,
    alamat TEXT NOT NULL,
    provinsi VARCHAR(100) DEFAULT '',
    kota VARCHAR(100) DEFAULT '',
    kecamatan VARCHAR(100) DEFAULT '',
    kode_pos VARCHAR(10) DEFAULT '',
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,

    -- Detail Lokasi
    luas_tempat DOUBLE PRECISION DEFAULT 0,
    harga_sewa_per_tahun BIGINT DEFAULT 0,
    durasi_sewa INT DEFAULT 1,
    tipe_bangunan VARCHAR(50) DEFAULT '',
    lebar_jalan DOUBLE PRECISION DEFAULT 0,
    jumlah_lantai INT DEFAULT 1,

    -- Traffic & Potensi
    estimasi_lalu_lintas INT DEFAULT 0,
    dekat_dengan TEXT DEFAULT '',
    jumlah_kompetitor INT DEFAULT 0,
    target_market TEXT DEFAULT '',

    -- Scoring
    score_traffic INT DEFAULT 0,
    score_sewa INT DEFAULT 0,
    score_kompetitor INT DEFAULT 0,
    score_akses INT DEFAULT 0,
    score_market INT DEFAULT 0,
    total_score INT DEFAULT 0,
    score_category VARCHAR(30) DEFAULT '',

    -- Status
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    catatan TEXT DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Location Surveys
CREATE TABLE IF NOT EXISTS location_surveys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES location_submissions(id) ON DELETE CASCADE,
    survey_date TIMESTAMPTZ,
    survey_by UUID REFERENCES users(id),
    hasil_survey TEXT DEFAULT '',
    catatan_survey TEXT DEFAULT '',
    estimasi_omzet BIGINT DEFAULT 0,
    estimasi_bep INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Location Files
CREATE TABLE IF NOT EXISTS location_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES location_submissions(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    file_type VARCHAR(50) NOT NULL DEFAULT 'photo',
    label VARCHAR(100) DEFAULT '',
    uploaded_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Location Approvals
CREATE TABLE IF NOT EXISTS location_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES location_submissions(id) ON DELETE CASCADE,
    approved_by UUID NOT NULL REFERENCES users(id),
    decision VARCHAR(30) NOT NULL DEFAULT 'approved',
    note TEXT DEFAULT '',
    approved_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_location_submissions_mitra ON location_submissions(mitra_id);
CREATE INDEX IF NOT EXISTS idx_location_submissions_status ON location_submissions(status);
CREATE INDEX IF NOT EXISTS idx_location_surveys_location ON location_surveys(location_id);
CREATE INDEX IF NOT EXISTS idx_location_files_location ON location_files(location_id);
CREATE INDEX IF NOT EXISTS idx_location_approvals_location ON location_approvals(location_id);
