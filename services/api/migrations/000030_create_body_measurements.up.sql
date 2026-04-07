CREATE TABLE body_measurements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    measured_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Anthropometric
    weight_kg       NUMERIC(6,2) CHECK (weight_kg IS NULL OR (weight_kg >= 1 AND weight_kg <= 500)),
    height_cm       NUMERIC(5,1) CHECK (height_cm IS NULL OR (height_cm >= 30 AND height_cm <= 300)),
    bmi             NUMERIC(5,2),           -- auto-calculated

    -- Body composition
    body_fat_pct    NUMERIC(5,2) CHECK (body_fat_pct IS NULL OR (body_fat_pct >= 0 AND body_fat_pct <= 100)),
    lean_mass_kg    NUMERIC(6,2),
    fat_mass_kg     NUMERIC(6,2),
    muscle_mass_kg  NUMERIC(6,2),
    bone_mass_kg    NUMERIC(5,2),
    water_pct       NUMERIC(5,2) CHECK (water_pct IS NULL OR (water_pct >= 0 AND water_pct <= 100)),
    visceral_fat    NUMERIC(5,1),
    basal_metabolic_rate INT,

    -- Circumferences (cm)
    waist_cm        NUMERIC(5,1),
    hip_cm          NUMERIC(5,1),
    chest_cm        NUMERIC(5,1),
    right_arm_cm    NUMERIC(5,1),
    left_arm_cm     NUMERIC(5,1),
    right_thigh_cm  NUMERIC(5,1),
    left_thigh_cm   NUMERIC(5,1),
    right_calf_cm   NUMERIC(5,1),
    left_calf_cm    NUMERIC(5,1),
    neck_cm         NUMERIC(5,1),
    abdomen_cm      NUMERIC(5,1),

    -- Skinfolds (mm)
    triceps_sf_mm   NUMERIC(5,1),
    biceps_sf_mm    NUMERIC(5,1),
    subscapular_sf_mm NUMERIC(5,1),
    suprailiac_sf_mm  NUMERIC(5,1),
    abdominal_sf_mm   NUMERIC(5,1),
    thigh_sf_mm       NUMERIC(5,1),
    calf_sf_mm        NUMERIC(5,1),

    source          TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual','device','import')),
    device_model    TEXT,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_body_measurements_tenant  ON body_measurements(tenant_id);
CREATE INDEX idx_body_measurements_patient ON body_measurements(patient_id);
CREATE INDEX idx_body_measurements_date    ON body_measurements(patient_id, measured_at);

CREATE TABLE measurement_publications (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    measurement_id  UUID NOT NULL REFERENCES body_measurements(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    published_by    UUID NOT NULL REFERENCES users(id),
    published_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    message         TEXT,

    UNIQUE (measurement_id)
);

CREATE INDEX idx_meas_pub_tenant  ON measurement_publications(tenant_id);
CREATE INDEX idx_meas_pub_patient ON measurement_publications(patient_id);
