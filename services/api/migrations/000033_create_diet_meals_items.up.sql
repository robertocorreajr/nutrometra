CREATE TABLE diet_meals (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    diet_id     UUID NOT NULL REFERENCES diets(id) ON DELETE CASCADE,
    meal_name   TEXT NOT NULL,
    meal_order  INT NOT NULL,
    notes       TEXT,
    UNIQUE (diet_id, meal_order)
);
CREATE INDEX idx_diet_meals_diet ON diet_meals(diet_id);

CREATE TABLE diet_meal_items (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES tenants(id),
    diet_meal_id         UUID NOT NULL REFERENCES diet_meals(id) ON DELETE CASCADE,
    food_item_id         UUID NOT NULL REFERENCES food_items(id),
    quantity_value       NUMERIC(8,2) NOT NULL,
    quantity_unit        TEXT NOT NULL,
    household_measure_id UUID REFERENCES household_measures(id),
    amount_description   TEXT,
    preparation_notes    TEXT,
    sort_order           INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_diet_meal_items_meal ON diet_meal_items(diet_meal_id);

CREATE TABLE diet_substitutions (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants(id),
    diet_meal_item_id       UUID NOT NULL REFERENCES diet_meal_items(id) ON DELETE CASCADE,
    substitute_food_item_id UUID NOT NULL REFERENCES food_items(id),
    quantity_value          NUMERIC(8,2),
    quantity_unit           TEXT,
    notes                   TEXT,
    sort_order              INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_diet_subs_item ON diet_substitutions(diet_meal_item_id);
