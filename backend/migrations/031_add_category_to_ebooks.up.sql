ALTER TABLE ebooks ADD COLUMN category VARCHAR(100) DEFAULT '' NOT NULL;
CREATE INDEX idx_ebooks_category ON ebooks(category);
