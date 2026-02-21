-- Migration: 022_add_title_to_agreements
-- Add title column to agreements table

ALTER TABLE agreements ADD COLUMN IF NOT EXISTS title VARCHAR(255) NOT NULL DEFAULT '';
