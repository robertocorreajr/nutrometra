CREATE TABLE roles (
    id                UUID PRIMARY KEY,
    code              VARCHAR(64) NOT NULL UNIQUE,
    application_scope VARCHAR(32) NOT NULL CHECK (application_scope IN ('tenant','backoffice')),
    name              VARCHAR(128) NOT NULL,
    description       TEXT
);

CREATE TABLE permissions (
    id                UUID PRIMARY KEY,
    code              VARCHAR(128) NOT NULL UNIQUE,
    application_scope VARCHAR(32) NOT NULL CHECK (application_scope IN ('tenant','backoffice')),
    description       TEXT
);

CREATE TABLE role_permissions (
    id            UUID PRIMARY KEY,
    role_id       UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    UNIQUE (role_id, permission_id)
);
