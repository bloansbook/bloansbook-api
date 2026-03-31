CREATE TABLE cash_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    txn_id TEXT NOT NULL UNIQUE,

    date DATE NOT NULL,

    type TEXT NOT NULL,
    amount NUMERIC(15,2) NOT NULL,

    account TEXT NOT NULL,

    linked_entity_type TEXT NOT NULL,
    linked_entity_id UUID NOT NULL,

    reference TEXT,

    created_by UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT cash_type_check CHECK (
        type IN ('inflow', 'outflow')
    ),

    CONSTRAINT cash_account_check CHECK (
        account IN ('cash', 'bank')
    ),

    CONSTRAINT cash_entity_type_check CHECK (
        linked_entity_type IN (
            'customer_payment',
            'bill_payment',
            'payroll_run',
            'adjustment'
        )
    ),

    CONSTRAINT cash_amount_check CHECK (
        amount > 0
    )
);