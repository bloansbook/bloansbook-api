CREATE TABLE IF NOT EXISTS system_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    setting_key TEXT NOT NULL UNIQUE,
    setting_value TEXT NOT NULL,

    data_type TEXT NOT NULL,
    description TEXT NOT NULL,
    module TEXT NOT NULL,

    updated_by UUID REFERENCES staff(id),
    updated_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT data_type_check
        CHECK (data_type IN ('integer','decimal','text','boolean','time'))
);

CREATE INDEX idx_settings_module ON system_settings(module);