CREATE TABLE bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    bill_id TEXT NOT NULL UNIQUE,

    supplier_id UUID NOT NULL REFERENCES suppliers(id),

    purchase_order_id UUID REFERENCES purchase_orders(id),

    is_inventory_bill BOOLEAN NOT NULL DEFAULT FALSE,

    category TEXT NOT NULL,

    description TEXT NOT NULL,

    amount NUMERIC(15,2) NOT NULL,
    amount_paid NUMERIC(15,2) NOT NULL DEFAULT 0,

    status TEXT NOT NULL DEFAULT 'unpaid',

    due_date DATE,

    reversed_by_bill_id UUID 
        REFERENCES bills(id),

    reverses_bill_id UUID 
        REFERENCES bills(id),

    reversal_reason TEXT,

    created_by UUID NOT NULL REFERENCES staff(id),
    approved_by UUID REFERENCES staff(id),

    approved_at TIMESTAMPTZ,
    posted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT bill_category_check CHECK (
        category IN (
            'raw_materials',
            'printing',
            'logistics',
            'utilities',
            'artisans',
            'other'
        )
    ),

    CONSTRAINT bill_status_check CHECK (
        status IN ('unpaid', 'part_paid', 'paid', 'reversed')
    ),

    CONSTRAINT bill_amount_check CHECK (
        amount >= 0 AND amount_paid >= 0
    ),

    CONSTRAINT bill_inventory_rule CHECK (
        (is_inventory_bill = TRUE AND purchase_order_id IS NOT NULL)
        OR
        (is_inventory_bill = FALSE AND purchase_order_id IS NULL)
    ),

    CONSTRAINT bill_reversal_check CHECK (
        (reverses_bill_id IS NOT NULL AND reversal_reason IS NOT NULL)
        OR
        (reverses_bill_id IS NULL)
    )
);

CREATE INDEX idx_bills_supplier_id ON bills(supplier_id);
CREATE INDEX idx_bills_po_id ON bills(purchase_order_id);
CREATE INDEX idx_bills_status ON bills(status);
CREATE INDEX idx_bills_inventory_flag ON bills(is_inventory_bill);
CREATE INDEX idx_bills_due_date ON bills(due_date);
CREATE UNIQUE INDEX one_reversal_per_bill ON bills(reverses_bill_id) WHERE reverses_bill_id IS NOT NULL;

CREATE TABLE bill_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id TEXT NOT NULL UNIQUE,

    bill_id UUID NOT NULL REFERENCES bills(id),

    supplier_id UUID NOT NULL REFERENCES suppliers(id),

    date DATE NOT NULL,

    amount NUMERIC(15,2) NOT NULL,

    method TEXT NOT NULL,

    reference TEXT,

    recorded_by UUID NOT NULL REFERENCES staff(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT payment_method_check CHECK (
        method IN ('cash', 'bank_transfer', 'pos')
    ),

    CONSTRAINT payment_amount_check CHECK (
        amount > 0
    )
);

CREATE INDEX idx_bill_payments_bill_id ON bill_payments(bill_id);
CREATE INDEX idx_bill_payments_supplier_id ON bill_payments(supplier_id);
CREATE INDEX idx_bill_payments_date ON bill_payments(date);