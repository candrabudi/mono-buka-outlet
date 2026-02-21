-- Rollback: 013_create_outlet_categories_table
DROP INDEX IF EXISTS idx_outlets_category_id;
ALTER TABLE outlets DROP COLUMN IF EXISTS category_id;
DROP TABLE IF EXISTS outlet_categories;
