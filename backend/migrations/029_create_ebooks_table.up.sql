CREATE TABLE IF NOT EXISTS ebooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    author VARCHAR(255) DEFAULT '',
    cover_url TEXT DEFAULT '',
    file_url TEXT DEFAULT '',
    price BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    total_sold INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ebooks_slug ON ebooks(slug);
CREATE INDEX idx_ebooks_is_active ON ebooks(is_active);
