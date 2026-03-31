CREATE TABLE IF NOT EXISTS audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    log_id TEXT NOT NULL UNIQUE,

    timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),

    user_id UUID NOT NULL REFERENCES staff(id),

    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,

    before_json JSONB,
    after_json JSONB,

    reason TEXT,

    ip_address TEXT,
    user_agent TEXT,

    CONSTRAINT audit_action_check
        CHECK (action IN (
            'CREATE','UPDATE','APPROVE','REVERSE',
            'LOGIN','LOGOUT','EXPORT','REJECT',
            'TERMINATE','OVERRIDE'
        ))
);

CREATE INDEX idx_audit_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_entity_type ON audit_log(entity_type);
CREATE INDEX idx_audit_entity_id ON audit_log(entity_id);
CREATE INDEX idx_audit_timestamp ON audit_log(timestamp);
CREATE INDEX idx_audit_action ON audit_log(action);