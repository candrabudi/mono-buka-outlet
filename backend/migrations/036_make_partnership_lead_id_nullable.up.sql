-- Migration: 036_make_partnership_lead_id_nullable
-- Allow partnerships to be created without a lead (e.g. from mitra applications)

ALTER TABLE partnerships ALTER COLUMN lead_id DROP NOT NULL;
