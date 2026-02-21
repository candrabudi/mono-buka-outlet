-- Add leader_id to partnerships
ALTER TABLE partnerships ADD COLUMN leader_id UUID REFERENCES users(id);
CREATE INDEX idx_partnerships_leader ON partnerships(leader_id);
