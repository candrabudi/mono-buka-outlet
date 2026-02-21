-- Migration: 013_create_outlet_categories_table

CREATE TABLE IF NOT EXISTS outlet_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    icon VARCHAR(10),
    description VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_outlet_categories_slug ON outlet_categories(slug);
CREATE INDEX idx_outlet_categories_is_active ON outlet_categories(is_active);

-- Seed default categories
INSERT INTO outlet_categories (id, name, slug, icon, description, is_active, sort_order) VALUES
    (uuid_generate_v4(), 'Franchise', 'franchise', '🏢', 'Model bisnis franchise dengan lisensi merek dagang', true, 1),
    (uuid_generate_v4(), 'Licensing', 'licensing', '📋', 'Model lisensi penggunaan merek dan sistem', true, 2),
    (uuid_generate_v4(), 'Kemitraan', 'kemitraan', '🤝', 'Model kemitraan usaha bersama', true, 3),
    (uuid_generate_v4(), 'BOT', 'bot', '🔄', 'Build-Operate-Transfer', true, 4);

-- Add category_id column to outlets
ALTER TABLE outlets ADD COLUMN category_id UUID REFERENCES outlet_categories(id);

-- Migrate existing category data
UPDATE outlets SET category_id = (SELECT id FROM outlet_categories WHERE slug = outlets.category LIMIT 1)
    WHERE category IS NOT NULL AND category != '';

-- Create index for the new column
CREATE INDEX idx_outlets_category_id ON outlets(category_id);
