-- Migration: 023_add_type_to_agreements
-- Add type column: CONTRACT or DOCUMENT

ALTER TABLE agreements ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'CONTRACT';
