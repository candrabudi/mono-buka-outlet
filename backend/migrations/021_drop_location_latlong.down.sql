DROP INDEX IF EXISTS idx_location_submissions_partnership;
ALTER TABLE location_submissions DROP COLUMN IF EXISTS partnership_id;
ALTER TABLE location_submissions ADD COLUMN IF NOT EXISTS latitude DOUBLE PRECISION DEFAULT 0;
ALTER TABLE location_submissions ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION DEFAULT 0;
