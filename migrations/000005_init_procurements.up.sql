CREATE TABLE purchase_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    po_id TEXT NOT NULL UNIQUE,

    supplier_id UUID NOT NULL REFERENCES suppliers(id),

    date DATE NOT NULL,
    expected_delivery_date DATE,

    status TEXT NOT NULL DEFAULT 'draft',

    quantity_tolerance NUMERIC(5,2) NOT NULL DEFAULT 5.00,

    notes TEXT,

    created_by UUID NOT NULL REFERENCES staff(id),
    approved_by UUID REFERENCES staff(id),
    approved_at TIMESTAMPTZ,

    cancelled_by UUID REFERENCES staff(id),
    cancellation_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- ✅ status constraint
    CONSTRAINT po_status_check CHECK (
        status IN (
            'draft',
            'submitted',
            'supplier_confirmed',
            'delivered',
            'closed',
            'cancelled'
        )
    ),

    CONSTRAINT tolerance_range CHECK (
        quantity_tolerance >= 0 AND quantity_tolerance <= 100
    ),

    CONSTRAINT cancellation_check CHECK (
        (status = 'cancelled' AND cancellation_reason IS NOT NULL)
        OR
        (status != 'cancelled')
    )
);

CREATE INDEX idx_po_supplier_id ON purchase_orders(supplier_id);
CREATE INDEX idx_po_status ON purchase_orders(status);
CREATE INDEX idx_po_date ON purchase_orders(date);

CREATE TABLE purchase_order_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    purchase_order_id UUID NOT NULL 
        REFERENCES purchase_orders(id) ON DELETE CASCADE,

    item_type TEXT NOT NULL,

    material_id UUID REFERENCES materials(id),
    variant_id UUID REFERENCES product_variants(id),

    description TEXT,

    quantity_ordered NUMERIC(15,2) NOT NULL,
    quantity_delivered NUMERIC(15,2) NOT NULL DEFAULT 0,

    unit_of_measure TEXT NOT NULL,

    unit_cost NUMERIC(15,2) NOT NULL,
    total_cost NUMERIC(15,2) NOT NULL,

    is_fully_delivered BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pol_item_type_check CHECK (
        item_type IN ('raw_material', 'finished_product')
    ),

    CONSTRAINT pol_one_item CHECK (
        (material_id IS NOT NULL AND variant_id IS NULL)
        OR
        (material_id IS NULL AND variant_id IS NOT NULL)
    ),

    CONSTRAINT pol_positive_values CHECK (
        quantity_ordered > 0 AND
        quantity_delivered >= 0 AND
        unit_cost >= 0 AND
        total_cost >= 0
    )
);

CREATE INDEX idx_pol_po_id ON purchase_order_lines(purchase_order_id);
CREATE INDEX idx_pol_material_id ON purchase_order_lines(material_id);
CREATE INDEX idx_pol_variant_id ON purchase_order_lines(variant_id);

CREATE TABLE po_delivery_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id),

    purchase_order_line_id UUID NOT NULL REFERENCES purchase_order_lines(id),

    delivery_date DATE NOT NULL,

    quantity_delivered NUMERIC(15,2) NOT NULL,

    notes TEXT,

    confirmed_by UUID NOT NULL REFERENCES staff(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pdl_positive_quantity CHECK (
        quantity_delivered > 0
    )
);

CREATE INDEX idx_pdl_po_id ON po_delivery_logs(purchase_order_id);
CREATE INDEX idx_pdl_line_id ON po_delivery_logs(purchase_order_line_id);
CREATE INDEX idx_pdl_date ON po_delivery_logs(delivery_date);