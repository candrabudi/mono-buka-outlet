-- Rollback: 040_add_referral_code_to_users

ALTER TABLE users DROP COLUMN IF EXISTS referral_code;
