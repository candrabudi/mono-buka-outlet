-- Migration: 043_create_affiliator_withdrawals
-- Withdrawal requests for affiliators (manual transfer by admin)

CREATE TABLE IF NOT EXISTS affiliator_withdrawals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    affiliator_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    amount DECIMAL(15,2) NOT NULL,
    bank_name VARCHAR(100) NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    account_holder VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    admin_notes TEXT,
    processed_by UUID REFERENCES users(id),
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_aff_withdrawals_affiliator ON affiliator_withdrawals(affiliator_id);
CREATE INDEX idx_aff_withdrawals_status ON affiliator_withdrawals(status);
CREATE INDEX idx_aff_withdrawals_created ON affiliator_withdrawals(created_at);
