-- Rollback: 036_make_partnership_lead_id_nullable
-- WARNING: This will fail if there are NULL lead_id values

ALTER TABLE partnerships ALTER COLUMN lead_id SET NOT NULL;
