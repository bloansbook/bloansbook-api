CREATE TABLE IF NOT EXISTS attendance_daily (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    staff_id UUID NOT NULL REFERENCES staff(id),

    date DATE NOT NULL,

    clock_in TIMESTAMPTZ,
    clock_out TIMESTAMPTZ,

    minutes_worked INTEGER NOT NULL DEFAULT 0,
    late_minutes INTEGER NOT NULL DEFAULT 0,

    observed_overtime_weekday INTEGER NOT NULL DEFAULT 0,
    observed_overtime_saturday INTEGER NOT NULL DEFAULT 0,

    is_absent BOOLEAN NOT NULL DEFAULT false,

    is_exception BOOLEAN NOT NULL DEFAULT false,
    exception_reason TEXT,

    exception_resolved BOOLEAN NOT NULL DEFAULT false,
    exception_resolved_by UUID REFERENCES staff(id),
    exception_resolved_at TIMESTAMPTZ,

    is_manual_edit BOOLEAN NOT NULL DEFAULT false,
    manual_edit_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT unique_staff_date UNIQUE (staff_id, date),

    CONSTRAINT exception_resolution_check
        CHECK (
            exception_resolved = false
            OR (exception_resolved_by IS NOT NULL AND exception_resolved_at IS NOT NULL)
        ),

    CONSTRAINT manual_edit_reason_check
        CHECK (
            is_manual_edit = false
            OR manual_edit_reason IS NOT NULL
        )
);

CREATE INDEX idx_attendance_staff_id ON attendance_daily(staff_id);
CREATE INDEX idx_attendance_date ON attendance_daily(date);
CREATE INDEX idx_attendance_is_exception ON attendance_daily(is_exception);
CREATE INDEX idx_attendance_is_absent ON attendance_daily(is_absent);

CREATE TABLE IF NOT EXISTS overtime_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    staff_id UUID NOT NULL REFERENCES staff(id),
    attendance_daily_id UUID NOT NULL REFERENCES attendance_daily(id),

    date DATE NOT NULL,

    observed_minutes_weekday INTEGER NOT NULL DEFAULT 0,
    observed_minutes_saturday INTEGER NOT NULL DEFAULT 0,

    approved_minutes_weekday INTEGER NOT NULL DEFAULT 0,
    approved_minutes_saturday INTEGER NOT NULL DEFAULT 0,

    status TEXT NOT NULL DEFAULT 'pending',

    reviewed_by UUID REFERENCES staff(id),
    reviewed_at TIMESTAMPTZ,
    review_notes TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT unique_date_staff UNIQUE (staff_id, date),

    CONSTRAINT overtime_status_check
        CHECK (status IN ('pending', 'approved', 'partial', 'rejected')),

    CONSTRAINT review_required_check
        CHECK (
            status = 'pending'
            OR (reviewed_by IS NOT NULL AND reviewed_at IS NOT NULL)
        ),

    CONSTRAINT review_notes_required_check
        CHECK (
            status NOT IN ('partial', 'rejected')
            OR review_notes IS NOT NULL
        )
);

CREATE INDEX idx_overtime_staff_id ON overtime_requests(staff_id);
CREATE INDEX idx_overtime_date ON overtime_requests(date);
CREATE INDEX idx_overtime_status ON overtime_requests(status);
CREATE INDEX idx_overtime_attendance_id ON overtime_requests(attendance_daily_id);