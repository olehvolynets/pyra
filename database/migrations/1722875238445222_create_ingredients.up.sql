CREATE TYPE ingredient_type AS ENUM ('none', 'product', 'dish');

CREATE TABLE IF NOT EXISTS ingredients
(
    dish_uid                UUID    NOT NULL,
    dish_version            int4    NOT NULL,

    ingredientable_uid      UUID    NOT NULL,
    ingredientable_version  int4    NOT NULL,
    ingredientable_type     ingredient_type NOT NULL,

    amount                  int4    NOT NULL,
    unit                    int2    NOT NULL,

    idx                     int2    NOT NULL,

    PRIMARY KEY (
        dish_uid,
        dish_version,
        ingredientable_uid,
        ingredientable_version,
        ingredientable_type
    )
);
