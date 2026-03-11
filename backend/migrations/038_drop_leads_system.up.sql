-- Migration: 038_drop_leads_system
-- Completely remove the leads system (table + FK references)

-- Drop lead_id FK from partnerships (already nullable via migration 036)
ALTER TABLE partnerships DROP COLUMN IF EXISTS lead_id;

-- Drop lead_id FK from locations
ALTER TABLE locations DROP COLUMN IF EXISTS lead_id;

-- Drop leads table entirely (CASCADE to handle any remaining deps)
DROP TABLE IF EXISTS leads CASCADE;
