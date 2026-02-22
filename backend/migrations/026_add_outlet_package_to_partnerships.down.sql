-- Rollback: 026_add_outlet_package_to_partnerships

DROP INDEX IF EXISTS idx_partnerships_package_id;
DROP INDEX IF EXISTS idx_partnerships_outlet_id;
ALTER TABLE partnerships DROP COLUMN IF EXISTS package_id;
ALTER TABLE partnerships DROP COLUMN IF EXISTS outlet_id;
