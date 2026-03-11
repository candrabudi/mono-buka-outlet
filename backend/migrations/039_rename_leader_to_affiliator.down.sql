-- Rollback: 039_rename_leader_to_affiliator

ALTER TABLE partnerships RENAME COLUMN affiliator_id TO leader_id;
UPDATE users SET role = 'leader' WHERE role = 'affiliator';
