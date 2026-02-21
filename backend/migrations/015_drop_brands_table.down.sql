-- Rollback: recreate brands table (if needed)
CREATE TABLE IF NOT EXISTS brands (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    logo VARCHAR(500),
    description TEXT,
    minimum_investment DECIMAL(15,2) NOT NULL DEFAULT 0,
    profit_sharing_percentage DECIMAL(5,2) NOT NULL DEFAULT 0,
    estimated_roi VARCHAR(100),
    location_requirement TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Re-add FK constraints
ALTER TABLE leads ADD CONSTRAINT leads_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE RESTRICT;
ALTER TABLE partnerships ADD CONSTRAINT partnerships_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE RESTRICT;

-- Re-add indexes
CREATE INDEX IF NOT EXISTS idx_leads_brand_id ON leads(brand_id);
CREATE INDEX IF NOT EXISTS idx_partnerships_brand_id ON partnerships(brand_id);
