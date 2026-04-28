CREATE TABLE job_labour (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL REFERENCES orders(id),

    staff_id UUID NOT NULL REFERENCES staff(id),

    hours NUMERIC(8,2) NOT NULL,
    rate NUMERIC(15,2) NOT NULL,
    total NUMERIC(15,2) NOT NULL,

    notes TEXT,

    recorded_by UUID NOT NULL REFERENCES staff(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT job_labor_positive CHECK (
        hours > 0 AND rate >= 0 AND total >= 0
    )
);

CREATE TABLE job_overhead (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL REFERENCES orders(id),

    type TEXT NOT NULL,

    amount NUMERIC(15,2) NOT NULL,

    basis TEXT,
    notes TEXT,

    recorded_by UUID NOT NULL REFERENCES staff(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT job_overhead_type_check CHECK (
        type IN ('generator', 'machine', 'miscellaneous')
    ),

    CONSTRAINT job_overhead_amount_check CHECK (
        amount >= 0
    )
);