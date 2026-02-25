-- Rollback: 037_make_partnership_brand_id_nullable

ALTER TABLE partnerships ALTER COLUMN brand_id SET NOT NULL;
