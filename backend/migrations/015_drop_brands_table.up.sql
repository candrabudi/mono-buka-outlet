-- Migration: 015_drop_brands_table
-- Remove foreign key constraints referencing brands table, then drop brands table

-- Drop FK constraints from leads
ALTER TABLE leads DROP CONSTRAINT IF EXISTS leads_brand_id_fkey;

-- Drop FK constraints from partnerships
ALTER TABLE partnerships DROP CONSTRAINT IF EXISTS partnerships_brand_id_fkey;

-- Drop FK constraints from payments (if any)
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_brand_id_fkey;

-- Drop FK constraints from revenues (if any)  
ALTER TABLE revenues DROP CONSTRAINT IF EXISTS revenues_brand_id_fkey;

-- Drop FK constraints from agreements (if any)
ALTER TABLE agreements DROP CONSTRAINT IF EXISTS agreements_brand_id_fkey;

-- Drop FK constraints from locations (if any)
ALTER TABLE locations DROP CONSTRAINT IF EXISTS locations_brand_id_fkey;

-- Drop brand index on leads
DROP INDEX IF EXISTS idx_leads_brand_id;

-- Drop brand index on partnerships
DROP INDEX IF EXISTS idx_partnerships_brand_id;

-- Finally drop the brands table
DROP TABLE IF EXISTS brands;
