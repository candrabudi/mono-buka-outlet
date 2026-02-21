-- Migration: 012_create_outlets_table

CREATE TABLE IF NOT EXISTS outlets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    logo VARCHAR(500),
    banner VARCHAR(500),
    category VARCHAR(100) NOT NULL DEFAULT 'franchise',
    description TEXT,
    short_description VARCHAR(500),
    minimum_investment DECIMAL(15,2) NOT NULL DEFAULT 0,
    maximum_investment DECIMAL(15,2),
    profit_sharing_percentage DECIMAL(5,2) NOT NULL DEFAULT 0,
    estimated_roi VARCHAR(100),
    location_requirement TEXT,
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    contact_phone VARCHAR(20),
    contact_email VARCHAR(255),
    contact_whatsapp VARCHAR(20),
    website VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_featured BOOLEAN NOT NULL DEFAULT false,
    total_outlets INTEGER NOT NULL DEFAULT 0,
    year_established INTEGER,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_outlets_slug ON outlets(slug);
CREATE INDEX idx_outlets_category ON outlets(category);
CREATE INDEX idx_outlets_is_active ON outlets(is_active);
CREATE INDEX idx_outlets_is_featured ON outlets(is_featured);
CREATE INDEX idx_outlets_city ON outlets(city);
CREATE INDEX idx_outlets_province ON outlets(province);
CREATE INDEX idx_outlets_deleted_at ON outlets(deleted_at);
