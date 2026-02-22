ALTER TABLE partnership_applications
  ADD COLUMN IF NOT EXISTS contact_phone VARCHAR(20) DEFAULT '',
  ADD COLUMN IF NOT EXISTS contact_email VARCHAR(255) DEFAULT '';
