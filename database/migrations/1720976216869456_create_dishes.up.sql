CREATE TABLE IF NOT EXISTS dishes
(
    uid         UUID                NOT NULL,
    version     int4                NOT NULL,
    name        VARCHAR(64)         NOT NULL,
    calories    int4                NOT NULL DEFAULT 0,
    proteins    int4                NOT NULL DEFAULT 0,
    fats        int4                NOT NULL DEFAULT 0,
    carbs       int4                NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(uid, version)
);
