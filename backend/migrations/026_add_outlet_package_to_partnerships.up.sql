-- Migration: 026_add_outlet_package_to_partnerships
-- Add outlet_id and package_id columns to partnerships table

ALTER TABLE partnerships ADD COLUMN IF NOT EXISTS outlet_id UUID REFERENCES outlets(id);
ALTER TABLE partnerships ADD COLUMN IF NOT EXISTS package_id UUID REFERENCES outlet_packages(id);

CREATE INDEX IF NOT EXISTS idx_partnerships_outlet_id ON partnerships(outlet_id);
CREATE INDEX IF NOT EXISTS idx_partnerships_package_id ON partnerships(package_id);
