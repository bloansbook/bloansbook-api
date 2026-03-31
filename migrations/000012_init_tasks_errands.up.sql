CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    task_id TEXT NOT NULL UNIQUE,

    task_date DATE NOT NULL,

    assigned_to UUID NOT NULL REFERENCES staff(id),
    assigned_by UUID NOT NULL REFERENCES staff(id),

    department TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,

    status TEXT NOT NULL DEFAULT 'assigned',

    is_offsite BOOLEAN NOT NULL DEFAULT false,

    due_date DATE NOT NULL,
    notes TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT task_status_check
        CHECK (status IN ('assigned', 'in_progress', 'done', 'not_done')),

    CONSTRAINT task_department_check
        CHECK (department IN ('factory', 'admin'))
);

CREATE INDEX idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX idx_tasks_assigned_by ON tasks(assigned_by);
CREATE INDEX idx_tasks_date ON tasks(task_date);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_department ON tasks(department);

CREATE TABLE IF NOT EXISTS checklist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,

    description TEXT NOT NULL,

    is_completed BOOLEAN NOT NULL DEFAULT false,
    completed_by UUID REFERENCES staff(id),
    completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT completion_check
        CHECK (
            is_completed = false
            OR (completed_by IS NOT NULL AND completed_at IS NOT NULL)
        )
);

CREATE INDEX idx_checklist_task_id ON checklist_items(task_id);

CREATE TABLE IF NOT EXISTS errands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    errand_id TEXT NOT NULL UNIQUE,

    task_id UUID REFERENCES tasks(id),

    staff_id UUID NOT NULL REFERENCES staff(id),

    department TEXT NOT NULL,
    purpose TEXT NOT NULL,
    destination TEXT NOT NULL,

    status TEXT NOT NULL DEFAULT 'requested',

    requested_by UUID NOT NULL REFERENCES staff(id),

    approved_by UUID REFERENCES staff(id),
    approved_at TIMESTAMPTZ,
    approval_notes TEXT,

    time_out TIMESTAMPTZ,
    time_returned TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT errand_status_check
        CHECK (status IN ('requested', 'approved', 'rejected', 'in_transit', 'completed')),

    CONSTRAINT errand_department_check
        CHECK (department IN ('factory', 'admin'))
);

CREATE INDEX idx_errands_task_id ON errands(task_id);
CREATE INDEX idx_errands_staff_id ON errands(staff_id);
CREATE INDEX idx_errands_status ON errands(status);
CREATE INDEX idx_errands_department ON errands(department);

CREATE TABLE IF NOT EXISTS performance_flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    staff_id UUID NOT NULL REFERENCES staff(id),

    month TEXT NOT NULL,

    incomplete_task_count INTEGER NOT NULL,

    status TEXT NOT NULL DEFAULT 'open',

    reviewed_by UUID REFERENCES staff(id),
    reviewed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT unique_staff_month UNIQUE (staff_id, month),

    CONSTRAINT performance_status_check
        CHECK (status IN ('open', 'reviewed'))
);

CREATE INDEX idx_perf_staff_id ON performance_flags(staff_id);
CREATE INDEX idx_perf_status ON performance_flags(status);
CREATE INDEX idx_perf_month ON performance_flags(month);