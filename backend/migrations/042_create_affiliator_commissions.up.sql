-- Migration: 042_create_affiliator_commissions
-- Commission log for affiliators

CREATE TABLE IF NOT EXISTS affiliator_commissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    affiliator_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    partnership_id UUID REFERENCES partnerships(id) ON DELETE SET NULL,
    amount DECIMAL(15,2) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'COMMISSION',
    description TEXT,
    given_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_aff_commissions_affiliator ON affiliator_commissions(affiliator_id);
CREATE INDEX idx_aff_commissions_partnership ON affiliator_commissions(partnership_id);
CREATE INDEX idx_aff_commissions_created ON affiliator_commissions(created_at);
