-- Add type & proof_url columns that the code already uses but migration was missing
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS type VARCHAR(20) NOT NULL DEFAULT 'INVOICE';
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS proof_url TEXT DEFAULT '';
