CREATE TABLE backoffice_users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL UNIQUE REFERENCES users(id),
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE backoffice_user_roles (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    backoffice_user_id UUID NOT NULL REFERENCES backoffice_users(id) ON DELETE CASCADE,
    role_id            UUID NOT NULL REFERENCES roles(id),
    UNIQUE (backoffice_user_id, role_id)
);
