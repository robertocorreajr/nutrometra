CREATE TABLE food_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID,                                  -- NULL = global/system item
    name            TEXT NOT NULL,
    food_group      TEXT NOT NULL,
    brand           TEXT,
    barcode         TEXT,
    serving_size_g  NUMERIC(8,2) NOT NULL,
    serving_label   TEXT NOT NULL DEFAULT '100g',
    source          TEXT NOT NULL DEFAULT 'system'
                    CHECK (source IN ('system','taco','ibge','tenant')),
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_food_items_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE INDEX idx_food_items_tenant ON food_items(tenant_id);
CREATE INDEX idx_food_items_group  ON food_items(food_group);
CREATE INDEX idx_food_items_name   ON food_items USING gin (to_tsvector('portuguese', name));

CREATE TABLE nutrition_facts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_item_id    UUID NOT NULL UNIQUE REFERENCES food_items(id),
    calories_kcal   NUMERIC(8,2),
    protein_g       NUMERIC(8,2),
    carbs_g         NUMERIC(8,2),
    fiber_g         NUMERIC(8,2),
    sugar_g         NUMERIC(8,2),
    total_fat_g     NUMERIC(8,2),
    saturated_fat_g NUMERIC(8,2),
    trans_fat_g     NUMERIC(8,2),
    cholesterol_mg  NUMERIC(8,2),
    sodium_mg       NUMERIC(8,2),
    potassium_mg    NUMERIC(8,2),
    calcium_mg      NUMERIC(8,2),
    iron_mg         NUMERIC(8,2),
    vitamin_a_mcg   NUMERIC(8,2),
    vitamin_c_mg    NUMERIC(8,2),
    vitamin_d_mcg   NUMERIC(8,2),
    vitamin_b12_mcg NUMERIC(8,2),
    zinc_mg         NUMERIC(8,2),
    magnesium_mg    NUMERIC(8,2)
);

CREATE TABLE household_measures (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_item_id    UUID NOT NULL REFERENCES food_items(id),
    label           TEXT NOT NULL,           -- ex: "colher de sopa", "xícara"
    grams           NUMERIC(8,2) NOT NULL,   -- peso em gramas
    UNIQUE (food_item_id, label)
);

CREATE INDEX idx_household_measures_food ON household_measures(food_item_id);
