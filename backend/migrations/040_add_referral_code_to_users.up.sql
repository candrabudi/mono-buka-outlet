-- Migration: 040_add_referral_code_to_users
-- Add referral_code column for affiliator referral system

ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_code VARCHAR(20) UNIQUE;

-- Generate referral codes for existing affiliators
UPDATE users
SET referral_code = UPPER(SUBSTRING(MD5(id::text || created_at::text) FROM 1 FOR 8))
WHERE role = 'affiliator' AND referral_code IS NULL;
