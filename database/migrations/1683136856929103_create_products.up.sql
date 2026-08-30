CREATE TABLE IF NOT EXISTS products
(
    uid             UUID        NOT NULL DEFAULT gen_random_uuid(),
    version         int4        NOT NULL DEFAULT 1,
    name            VARCHAR(64) NOT NULL,
    calories        int4        NOT NULL DEFAULT 0,
    proteins        int4        NOT NULL DEFAULT 0,
    fats            int4        NOT NULL DEFAULT 0,
    carbs           int4        NOT NULL DEFAULT 0,
    archived_at     TIMESTAMPTZ, 
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (uid, version)
);

CREATE INDEX idx_product_on_archived ON products (archived_at);
