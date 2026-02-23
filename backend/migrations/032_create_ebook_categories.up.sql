-- Drop the simple category column (replaced by M2M relationship)
ALTER TABLE ebooks DROP COLUMN IF EXISTS category;
DROP INDEX IF EXISTS idx_ebooks_category;

-- Master table for ebook categories
CREATE TABLE IF NOT EXISTS ebook_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ebook_categories_slug ON ebook_categories(slug);

-- Many-to-many pivot table
CREATE TABLE IF NOT EXISTS ebook_category_mapping (
    ebook_id UUID NOT NULL REFERENCES ebooks(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES ebook_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (ebook_id, category_id)
);

CREATE INDEX idx_ecm_ebook ON ebook_category_mapping(ebook_id);
CREATE INDEX idx_ecm_category ON ebook_category_mapping(category_id);

-- Seed default categories
INSERT INTO ebook_categories (id, name, slug) VALUES
    (uuid_generate_v4(), 'Bisnis', 'bisnis'),
    (uuid_generate_v4(), 'Marketing', 'marketing'),
    (uuid_generate_v4(), 'Keuangan', 'keuangan'),
    (uuid_generate_v4(), 'Operasional', 'operasional'),
    (uuid_generate_v4(), 'Motivasi', 'motivasi');
