package dishes

const listDishesQuery = `
SELECT dishes.*
FROM dishes
INNER JOIN (
	SELECT DISTINCT uid, max(version) AS version
	FROM dishes
	GROUP BY uid
) latest_dishes ON dishes.uid = latest_dishes.uid
AND dishes.version = latest_dishes.version
ORDER BY created_at DESC;`

const dishByRefQuery = "SELECT * FROM dishes WHERE uid = $1 AND version = $2 LIMIT 1;"

const dishesByProductQuery = `
SELECT dishes.*
FROM dishes
INNER JOIN ingredients
	ON ingredients.dish_uid = dishes.uid
	AND ingredients.dish_version = dishes.version
WHERE
	ingredients.ingredientable_type = 'product'
	AND ingredients.ingredientable_uid = $1
	AND ingredients.ingredientable_version = $2;`

const dishVersionsQuery = `
SELECT *
FROM dishes
WHERE uid = $1
ORDER BY version DESC
LIMIT 20;`

const createDishQuery = `
INSERT INTO dishes (
	uid, version, name, calories, proteins, fats, carbs
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
) RETURNING created_at, updated_at;`

const isDishNameTakenQuery = `SELECT 1 FROM dishes WHERE name = $1 AND uid != $2;`

const deleteByRefQuery = `
DELETE FROM dishes
WHERE uid = $1 AND version = $2;`
// const createIngredientsQuery = `
// INSERT INTO ingredients (
// 	dish_id, ingredientable_id, ingredientable_type, amount, unit, idx
// ) VALUES (
// 	$1
// )`
