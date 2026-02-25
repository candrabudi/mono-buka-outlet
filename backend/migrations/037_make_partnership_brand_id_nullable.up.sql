-- Migration: 037_make_partnership_brand_id_nullable
-- Allow partnerships to be created without a brand

ALTER TABLE partnerships ALTER COLUMN brand_id DROP NOT NULL;
