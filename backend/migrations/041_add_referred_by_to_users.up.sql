-- Migration: 041_add_referred_by_to_users
-- Add referred_by column to track which affiliator referred a mitra

ALTER TABLE users ADD COLUMN IF NOT EXISTS referred_by UUID REFERENCES users(id);
CREATE INDEX IF NOT EXISTS idx_users_referred_by ON users(referred_by);
