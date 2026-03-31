CREATE TABLE IF NOT EXISTS customers (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),

    customer_id TEXT UNIQUE NOT NULL,

    name TEXT NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email TEXT,
    address TEXT,

    notes TEXT,
    type TEXT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'NGN',

    created_by UUID NOT NULL REFERENCES staff(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT customer_type_check CHECK (type IN ('retail', 'corporate'))
);

CREATE INDEX idx_customers_id ON customers(id);
CREATE INDEX idx_customers_customer_id ON customers(customer_id);
CREATE INDEX idx_customers_name ON customers(name);