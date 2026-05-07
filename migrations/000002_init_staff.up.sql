CREATE TABLE IF NOT EXISTS staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    staff_id TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT,
    phone TEXT,
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
    fired_at TIMESTAMPTZ,

    has_login boolean NOT NULL DEFAULT FALSE,
    supabase_uid TEXT UNIQUE,

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fired_requires_timestamp
    CHECK (
        status = 'fired' AND fired_at IS NOT NULL
        OR status != 'fired'
    ),

    CONSTRAINT valid_status CHECK (status IN ('active', 'inactive', 'fired'))
);
CREATE INDEX idx_staff_login ON staff(has_login, staff_id, status);
CREATE INDEX idx_staff_supabase_uid ON staff(supabase_uid);
CREATE INDEX idx_staff_active_login ON staff(has_login, staff_id) WHERE status = 'active';

CREATE TABLE IF NOT EXISTS fired_staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    staff_id UUID NOT NULL,

    termination_reason TEXT NOT NULL,

    recorded_by UUID NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    is_overridden BOOLEAN NOT NULL DEFAULT FALSE,
    overridden_by UUID,
    overridden_at TIMESTAMPTZ,
    override_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT override_requires_reason
    CHECK (
        is_overridden = FALSE
        OR override_reason IS NOT NULL
    )
);
CREATE INDEX idx_fired_staff_staff_id ON fired_staff(staff_id);

CREATE TABLE IF NOT EXISTS staff_roles (
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,

    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by UUID NOT NULL REFERENCES staff(id),

    PRIMARY KEY (staff_id, role_id)
);

CREATE INDEX idx_staff_roles_staff_id ON staff_roles(staff_id);
CREATE INDEX idx_staff_roles_role_id ON staff_roles(role_id);
CREATE INDEX idx_staff_roles_assigned_by ON staff_roles(assigned_by);

CREATE TABLE IF NOT EXISTS staff_role_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    staff_id UUID NOT NULL REFERENCES staff(id),
    role_id UUID NOT NULL REFERENCES roles(id),

    action TEXT NOT NULL,

    performed_by UUID NOT NULL,
    reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_actions CHECK (action IN ('assigned', 'revoked'))
);

CREATE INDEX idx_staff_role_history_staff_id ON staff_role_history(staff_id);
CREATE INDEX idx_staff_role_history_role_id ON staff_role_history(role_id);
CREATE INDEX idx_staff_role_history_performed_by ON staff_role_history(performed_by);