-- Refresh sessions: revocation via revoked_at (row kept for audit / replay detection).

CREATE TABLE refresh_sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash varchar(128) NOT NULL,
    user_agent text,
    ip_address varchar(45),
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    revoked_at timestamptz
);

CREATE UNIQUE INDEX uq_refresh_sessions_token_hash ON refresh_sessions (token_hash);

CREATE INDEX idx_refresh_sessions_user_id ON refresh_sessions (user_id);
CREATE INDEX idx_refresh_sessions_user_active ON refresh_sessions (user_id, expires_at)
WHERE revoked_at IS NULL;

CREATE INDEX idx_refresh_sessions_expires_at ON refresh_sessions (expires_at)
WHERE revoked_at IS NULL;
