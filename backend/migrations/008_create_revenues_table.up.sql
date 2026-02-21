-- Migration: 008_create_revenues_table

CREATE TABLE IF NOT EXISTS revenues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    partnership_id UUID NOT NULL REFERENCES partnerships(id) ON DELETE RESTRICT,
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE RESTRICT,
    month VARCHAR(7) NOT NULL,
    revenue DECIMAL(15,2) NOT NULL DEFAULT 0,
    expense DECIMAL(15,2) NOT NULL DEFAULT 0,
    profit DECIMAL(15,2) NOT NULL DEFAULT 0,
    company_share DECIMAL(15,2) NOT NULL DEFAULT 0,
    mitra_share DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_revenues_partnership_month ON revenues(partnership_id, month);
CREATE INDEX idx_revenues_brand_id ON revenues(brand_id);
CREATE INDEX idx_revenues_month ON revenues(month);
