-- Migration: 022_add_title_to_agreements (down)

ALTER TABLE agreements DROP COLUMN IF EXISTS title;
