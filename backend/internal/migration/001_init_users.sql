-- Users: account lifecycle via is_active; no deleted_at (hard delete policy at app layer if needed).

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email varchar(320) NOT NULL,
    password_hash varchar(255) NOT NULL,
    nickname varchar(64) NOT NULL,
    first_name varchar(100),
    last_name varchar(100),
    avatar_url text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    last_seen_at timestamptz,
    is_active boolean NOT NULL DEFAULT true
);

CREATE UNIQUE INDEX uq_users_email_lower ON users (lower(email));

CREATE INDEX idx_users_created_at ON users (created_at);
CREATE INDEX idx_users_last_seen_at ON users (last_seen_at) WHERE last_seen_at IS NOT NULL;
CREATE INDEX idx_users_is_active ON users (is_active) WHERE is_active = true;
