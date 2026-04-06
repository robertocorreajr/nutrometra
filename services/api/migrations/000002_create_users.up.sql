CREATE TABLE users (
    id                UUID PRIMARY KEY,
    email             VARCHAR(320) NOT NULL UNIQUE,
    email_verified_at TIMESTAMPTZ,
    password_hash     TEXT,
    external_auth_id  VARCHAR(255) UNIQUE,
    status            VARCHAR(32) NOT NULL DEFAULT 'active' CHECK (status IN ('active','disabled')),
    last_login_at     TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email            ON users(email);
CREATE INDEX idx_users_external_auth_id ON users(external_auth_id);
