-- Migration: 039_rename_leader_to_affiliator
-- Rename role 'leader' to 'affiliator' and column leader_id to affiliator_id

-- Rename column in partnerships
ALTER TABLE partnerships RENAME COLUMN leader_id TO affiliator_id;

-- Update role in users table
UPDATE users SET role = 'affiliator' WHERE role = 'leader';
