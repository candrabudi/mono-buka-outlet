ALTER TABLE location_submissions DROP COLUMN IF EXISTS latitude;
ALTER TABLE location_submissions DROP COLUMN IF EXISTS longitude;
ALTER TABLE location_submissions ADD COLUMN IF NOT EXISTS partnership_id UUID REFERENCES partnerships(id);
CREATE INDEX IF NOT EXISTS idx_location_submissions_partnership ON location_submissions(partnership_id);
