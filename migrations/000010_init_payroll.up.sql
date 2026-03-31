CREATE TABLE payroll_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    run_id TEXT NOT NULL UNIQUE,

    month TEXT NOT NULL UNIQUE,

    pay_date DATE,

    status TEXT NOT NULL DEFAULT 'draft',

    total_net_pay NUMERIC(15,2) NOT NULL DEFAULT 0,

    created_by UUID NOT NULL REFERENCES staff(id),
    approved_by UUID REFERENCES staff(id),

    approved_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT payroll_status_check CHECK (
        status IN ('draft', 'approved', 'paid')
    )
);

CREATE TABLE payroll_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    run_id UUID NOT NULL 
        REFERENCES payroll_runs(id) ON DELETE CASCADE,

    staff_id UUID NOT NULL 
        REFERENCES staff(id),

    base_salary NUMERIC(15,2) NOT NULL,

    overtime_hours_weekday NUMERIC(8,2) NOT NULL DEFAULT 0,
    overtime_hours_saturday NUMERIC(8,2) NOT NULL DEFAULT 0,

    overtime_rate_weekday NUMERIC(15,2) NOT NULL DEFAULT 0,
    overtime_rate_saturday NUMERIC(15,2) NOT NULL DEFAULT 0,

    overtime_earnings NUMERIC(15,2) NOT NULL DEFAULT 0,

    night_shift_count INTEGER NOT NULL DEFAULT 0,
    sunday_shift_count INTEGER NOT NULL DEFAULT 0,

    night_shift_earnings NUMERIC(15,2) NOT NULL DEFAULT 0,
    sunday_shift_earnings NUMERIC(15,2) NOT NULL DEFAULT 0,

    allowances NUMERIC(15,2) NOT NULL DEFAULT 0,
    allowances_notes TEXT,

    lateness_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,
    absence_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,

    other_deductions NUMERIC(15,2) NOT NULL DEFAULT 0,
    other_deductions_notes TEXT,

    net_pay NUMERIC(15,2) NOT NULL,

    is_manually_adjusted BOOLEAN NOT NULL DEFAULT FALSE,
    adjustment_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT payroll_line_unique UNIQUE (run_id, staff_id),

    CONSTRAINT payroll_line_adjustment_check CHECK (
        (is_manually_adjusted = TRUE AND adjustment_reason IS NOT NULL)
        OR
        (is_manually_adjusted = FALSE)
    )
);