CREATE TABLE IF NOT EXISTS ebook_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ebook_id UUID NOT NULL REFERENCES ebooks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number VARCHAR(50) NOT NULL UNIQUE,
    amount BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',

    -- Download approval workflow
    download_status VARCHAR(30) NOT NULL DEFAULT 'NONE',
    download_requested_at TIMESTAMPTZ,
    download_approved_at TIMESTAMPTZ,
    download_note TEXT DEFAULT '',

    -- Midtrans
    midtrans_order_id VARCHAR(100) DEFAULT '',
    midtrans_snap_token TEXT DEFAULT '',
    midtrans_redirect_url TEXT DEFAULT '',
    midtrans_payment_type VARCHAR(50) DEFAULT '',
    midtrans_transaction_id VARCHAR(100) DEFAULT '',
    midtrans_transaction_status VARCHAR(50) DEFAULT '',

    paid_at TIMESTAMPTZ,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ebook_orders_ebook_id ON ebook_orders(ebook_id);
CREATE INDEX idx_ebook_orders_user_id ON ebook_orders(user_id);
CREATE INDEX idx_ebook_orders_status ON ebook_orders(status);
CREATE INDEX idx_ebook_orders_download_status ON ebook_orders(download_status);
CREATE INDEX idx_ebook_orders_midtrans_order_id ON ebook_orders(midtrans_order_id);
