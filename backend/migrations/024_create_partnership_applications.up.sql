-- Migration: 024_create_partnership_applications

CREATE TABLE IF NOT EXISTS partnership_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mitra_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    outlet_id UUID NOT NULL REFERENCES outlets(id) ON DELETE RESTRICT,
    package_id UUID NOT NULL REFERENCES outlet_packages(id) ON DELETE RESTRICT,
    motivation TEXT,
    experience TEXT,
    proposed_location TEXT,
    investment_budget DECIMAL(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    admin_notes TEXT,
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_partnership_apps_mitra_id ON partnership_applications(mitra_id);
CREATE INDEX idx_partnership_apps_outlet_id ON partnership_applications(outlet_id);
CREATE INDEX idx_partnership_apps_package_id ON partnership_applications(package_id);
CREATE INDEX idx_partnership_apps_status ON partnership_applications(status);
