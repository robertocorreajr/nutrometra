CREATE TABLE coupons (
    id              UUID PRIMARY KEY,
    code            VARCHAR(64) NOT NULL UNIQUE,
    discount_type   VARCHAR(32) NOT NULL CHECK (discount_type IN ('percent','fixed')),
    discount_value  BIGINT NOT NULL,
    duration_type   VARCHAR(32) NOT NULL CHECK (duration_type IN ('once','repeating','forever')),
    duration_cycles INTEGER,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    starts_at       TIMESTAMPTZ,
    ends_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
