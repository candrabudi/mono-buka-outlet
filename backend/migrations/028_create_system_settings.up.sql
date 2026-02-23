CREATE TABLE IF NOT EXISTS system_settings (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  key VARCHAR(100) UNIQUE NOT NULL,
  value TEXT NOT NULL DEFAULT '',
  group_name VARCHAR(50) NOT NULL DEFAULT 'general',
  label VARCHAR(100) NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_system_settings_key ON system_settings(key);
CREATE INDEX idx_system_settings_group ON system_settings(group_name);

-- Seed Midtrans default settings
INSERT INTO system_settings (id, key, value, group_name, label, description, created_at, updated_at) VALUES
  (uuid_generate_v4(), 'midtrans_server_key', 'SB-Mid-server-XXXXXXXXXXXXXXXXXXXXXXXX', 'midtrans', 'Server Key', 'Midtrans Server Key (dari dashboard Midtrans)', NOW(), NOW()),
  (uuid_generate_v4(), 'midtrans_client_key', 'SB-Mid-client-XXXXXXXXXXXXXXXXXXXXXXXX', 'midtrans', 'Client Key', 'Midtrans Client Key (untuk Snap.js di frontend)', NOW(), NOW()),
  (uuid_generate_v4(), 'midtrans_environment', 'sandbox', 'midtrans', 'Environment', 'Environment Midtrans: sandbox atau production', NOW(), NOW())
ON CONFLICT (key) DO NOTHING;
