CREATE TABLE IF NOT EXISTS invoices (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  partnership_id UUID NOT NULL REFERENCES partnerships(id),
  invoice_number VARCHAR(50) UNIQUE NOT NULL,
  amount BIGINT NOT NULL DEFAULT 0,
  description TEXT DEFAULT '',
  status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
  
  -- Midtrans
  midtrans_order_id VARCHAR(100) UNIQUE,
  midtrans_snap_token TEXT,
  midtrans_redirect_url TEXT,
  midtrans_payment_type VARCHAR(50),
  midtrans_transaction_id VARCHAR(100),
  midtrans_transaction_status VARCHAR(50),
  
  paid_at TIMESTAMPTZ,
  expired_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoices_partnership ON invoices(partnership_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_order_id ON invoices(midtrans_order_id);
