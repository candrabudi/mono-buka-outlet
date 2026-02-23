ALTER TABLE ebook_orders ADD COLUMN IF NOT EXISTS payment_proof_url TEXT DEFAULT '';
