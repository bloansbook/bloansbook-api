CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id TEXT NOT NULL UNIQUE,

    customer_id UUID NOT NULL 
        REFERENCES customers(id),

    date DATE NOT NULL,
    expected_delivery_date DATE,

    status TEXT NOT NULL DEFAULT 'draft',

    customisation_notes TEXT,

    created_by UUID NOT NULL REFERENCES staff(id),
    approved_by UUID REFERENCES staff(id),
    approved_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT order_status_check CHECK (
        status IN (
            'draft',
            'approved',
            'in_production',
            'delivered',
            'closed'
        )
    )
);

CREATE TABLE order_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL 
        REFERENCES orders(id) ON DELETE CASCADE,

    variant_id UUID NOT NULL 
        REFERENCES product_variants(id),

    quantity NUMERIC(15,2) NOT NULL,
    quantity_delivered NUMERIC(15,2) NOT NULL DEFAULT 0,

    delivery_status TEXT NOT NULL DEFAULT 'pending',

    unit_price NUMERIC(15,2) NOT NULL,
    total_price NUMERIC(15,2) NOT NULL,

    notes TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT order_line_delivery_status_check CHECK (
        delivery_status IN ('pending', 'partial', 'delivered')
    ),

    CONSTRAINT order_line_positive_values CHECK (
        quantity > 0 AND
        quantity_delivered >= 0 AND
        unit_price >= 0 AND
        total_price >= 0
    )
);

CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    invoice_id TEXT NOT NULL UNIQUE,

    order_id UUID NOT NULL UNIQUE 
        REFERENCES orders(id),

    customer_id UUID NOT NULL 
        REFERENCES customers(id),

    date DATE NOT NULL,
    due_date DATE,

    subtotal NUMERIC(15,2) NOT NULL,

    discount_type TEXT,
    discount_value NUMERIC(15,2),
    discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0,

    total NUMERIC(15,2) NOT NULL,

    amount_paid NUMERIC(15,2) NOT NULL DEFAULT 0,

    status TEXT NOT NULL DEFAULT 'unpaid',

    discount_approval_status TEXT,

    discount_approved_by UUID REFERENCES staff(id),
    discount_approved_at TIMESTAMPTZ,

    reversed_by_invoice_id UUID REFERENCES invoices(id),
    reverses_invoice_id UUID REFERENCES invoices(id),

    reversal_reason TEXT,

    created_by UUID NOT NULL REFERENCES staff(id),

    posted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT invoice_status_check CHECK (
        status IN ('unpaid', 'part_paid', 'paid', 'cancelled', 'reversed')
    ),

    CONSTRAINT discount_type_check CHECK (
        discount_type IN ('percentage', 'fixed') OR discount_type IS NULL
    ),

    CONSTRAINT discount_approval_check CHECK (
        discount_approval_status IN ('pending', 'approved', 'rejected')
        OR discount_approval_status IS NULL
    ),

    CONSTRAINT invoice_amounts_check CHECK (
        subtotal >= 0 AND
        discount_amount >= 0 AND
        total >= 0 AND
        amount_paid >= 0
    ),

    CONSTRAINT invoice_reversal_check CHECK (
        (reverses_invoice_id IS NOT NULL AND reversal_reason IS NOT NULL)
        OR reverses_invoice_id IS NULL
    )
);

CREATE TABLE customer_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id TEXT NOT NULL UNIQUE,

    invoice_id UUID NOT NULL 
        REFERENCES invoices(id),

    customer_id UUID NOT NULL 
        REFERENCES customers(id),

    order_id UUID NOT NULL 
        REFERENCES orders(id),

    date DATE NOT NULL,

    amount NUMERIC(15,2) NOT NULL,

    method TEXT NOT NULL,

    reference TEXT,

    recorded_by UUID NOT NULL 
        REFERENCES staff(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT customer_payment_method_check CHECK (
        method IN ('cash', 'bank_transfer', 'pos')
    ),

    CONSTRAINT customer_payment_amount_check CHECK (
        amount > 0
    )
);