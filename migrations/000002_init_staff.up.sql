CREATE TABLE IF NOT EXISTS staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT,
    address TEXT,
    date_of_birth DATE,
    date_of_hire DATE,
    emergency_contact_name TEXT,
    emergency_contact_phone TEXT,
    bank_name TEXT,
    bank_account_number TEXT,
    bank_account_name TEXT,
    department TEXT NOT NULL,
    job_title TEXT NOT NULL,
    pay_type TEXT DEFAULT 'monthly' NOT NULL,
    base_salary NUMERIC(15, 2) NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    has_login boolean NOT NULL DEFAULT FALSE,
    superbase_uid TEXT UNIQUE,
    created_by TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_staff_login ON staff(has_login, staff_id, status);

CREATE TABLE IF NOT EXISTS fired_staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id TEXT NOT NULL,
    termination_date DATE NOT NULL,
    termination_reason TEXT NOT NULL,
    recorded_by TEXT NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_overridden BOOLEAN NOT NULL DEFAULT FALSE,
    overridden_by TEXT,
    overridden_at TIMESTAMPTZ,
    override_reason TEXT,

    CONSTRAINT override_requires_reason
    CHECK (
        is_overridden = false
        OR override_reason IS NOT NULL
    )
);