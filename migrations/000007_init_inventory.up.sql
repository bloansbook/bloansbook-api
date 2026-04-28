CREATE TABLE inventory_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    txn_id TEXT NOT NULL UNIQUE,

    idempotency_key TEXT,

    date DATE NOT NULL,

    type TEXT NOT NULL,
    item_type TEXT NOT NULL,

    material_id UUID REFERENCES materials(id),
    variant_id UUID REFERENCES product_variants(id),

    quantity NUMERIC(15,2) NOT NULL,

    unit_cost NUMERIC(15,2),
    total_cost NUMERIC(15,2),

    wac_at_transaction NUMERIC(15,2),

    order_id UUID REFERENCES orders(id),
    supplier_id UUID REFERENCES suppliers(id),
    purchase_order_id UUID REFERENCES purchase_orders(id),
    bill_id UUID REFERENCES bills(id),

    reversed_by UUID REFERENCES inventory_transactions(id),
    reverses UUID REFERENCES inventory_transactions(id),

    reversal_reason TEXT,

    notes TEXT,

    status TEXT NOT NULL DEFAULT 'approved',

    created_by UUID NOT NULL REFERENCES staff(id),
    approved_by UUID REFERENCES staff(id),
    approved_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT txn_type_check CHECK (
        type IN (
            'PURCHASE',
            'ISSUE_TO_JOB',
            'ADJUSTMENT_IN',
            'ADJUSTMENT_OUT',
            'PRODUCTION_IN',
            'SALE_OUT',
            'REVERSAL'
        )
    ),

    CONSTRAINT item_type_check CHECK (
        item_type IN ('raw_material', 'finished_product')
    ),

    CONSTRAINT status_check CHECK (
        status IN ('pending', 'approved', 'reversed')
    ),

    CONSTRAINT one_item_only CHECK (
        (material_id IS NOT NULL AND variant_id IS NULL)
        OR
        (material_id IS NULL AND variant_id IS NOT NULL)
    ),

    CONSTRAINT positive_quantity CHECK (quantity > 0),

    CONSTRAINT non_negative_costs CHECK (
        unit_cost IS NULL OR unit_cost >= 0
    ),

    CONSTRAINT total_cost_check CHECK (
        total_cost IS NULL OR total_cost >= 0
    ),

    CONSTRAINT reversal_consistency CHECK (
        (type = 'REVERSAL' AND reverses IS NOT NULL AND reversal_reason IS NOT NULL)
        OR
        (type != 'REVERSAL')
    )
);

CREATE UNIQUE INDEX idx_txn_idempotency ON inventory_transactions(idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE INDEX idx_txn_material_id ON inventory_transactions(material_id);
CREATE INDEX idx_txn_variant_id ON inventory_transactions(variant_id);
CREATE INDEX idx_txn_date ON inventory_transactions(date);
CREATE INDEX idx_txn_type ON inventory_transactions(type);
CREATE INDEX idx_txn_status ON inventory_transactions(status);
CREATE INDEX idx_txn_order_id ON inventory_transactions(order_id);
CREATE INDEX idx_txn_bill_id ON inventory_transactions(bill_id);
CREATE UNIQUE INDEX one_reversal_per_txn ON inventory_transactions(reverses) WHERE reverses IS NOT NULL;

CREATE TABLE wac_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    item_type TEXT NOT NULL,

    material_id UUID REFERENCES materials(id),
    variant_id UUID REFERENCES product_variants(id),

    wac NUMERIC(15,2) NOT NULL,
    previous_wac NUMERIC(15,2) NOT NULL,

    quantity_at_snapshot NUMERIC(15,2) NOT NULL,
    stock_value_at_snapshot NUMERIC(15,2) NOT NULL,

    triggered_by_txn_id UUID NOT NULL 
        REFERENCES inventory_transactions(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT snapshot_item_type_check CHECK (
        item_type IN ('raw_material', 'finished_product')
    ),

    CONSTRAINT snapshot_one_item CHECK (
        (material_id IS NOT NULL AND variant_id IS NULL)
        OR
        (material_id IS NULL AND variant_id IS NOT NULL)
    ),

    CONSTRAINT non_negative_snapshot CHECK (
        wac >= 0 AND
        previous_wac >= 0 AND
        quantity_at_snapshot >= 0 AND
        stock_value_at_snapshot >= 0
    )
);

CREATE INDEX idx_wac_material ON wac_snapshots(material_id);
CREATE INDEX idx_wac_variant ON wac_snapshots(variant_id);
CREATE INDEX idx_wac_created_at ON wac_snapshots(created_at);
CREATE INDEX idx_wac_txn ON wac_snapshots(triggered_by_txn_id);