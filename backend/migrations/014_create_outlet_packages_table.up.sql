CREATE TABLE IF NOT EXISTS outlet_packages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id UUID NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL DEFAULT 0,
    duration VARCHAR(100),
    image VARCHAR(500),
    estimated_bep VARCHAR(100),
    net_profit VARCHAR(100),
    description TEXT,
    benefits TEXT[],
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(outlet_id, slug)
);

CREATE INDEX idx_outlet_packages_outlet_id ON outlet_packages(outlet_id);
CREATE INDEX idx_outlet_packages_active ON outlet_packages(is_active);
