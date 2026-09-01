package products

const productColumns = "uid, version, name, calories, proteins, fats, carbs, created_at, updated_at"

const indexProductsQuery /* sql */ = `SELECT 
		products.uid,
		products.version,
		products.name,
		products.calories,
		products.proteins,
		products.fats,
		products.carbs,
		products.created_at,
		products.updated_at
	FROM products
	INNER JOIN (
		SELECT DISTINCT uid, max(version) AS version
		FROM products
		GROUP BY uid
	) latest_products ON products.uid = latest_products.uid
	AND products.version = latest_products.version
	ORDER BY created_at DESC;`

const findByRefQuery /* sql */ = "SELECT " + productColumns + ` FROM products
WHERE uid = $1 AND version = $2
LIMIT 1;`

const productVersionsQuery /* sql */ = "SELECT " + productColumns + ` FROM products
	WHERE uid = $1
	ORDER BY version DESC
	LIMIT 20;`

const findAllByIDsQuery /* sql */ = "SELECT " + productColumns + " FROM products WHERE id in $0;"

const productsForDishQuery /* sql */ = "SELECT " + productColumns + ` FROM ingredients
JOIN products
ON ingredients.ingredientable_uid = products.uid
	AND ingredients.ingredientable_version = products.version
	AND ingredients.ingredientable_type = 'product'
WHERE ingredients.dish_uid = $1 AND ingredients.dish_version = $2;`

const createProductVersionQuery /* sql */ = `
INSERT INTO products (
	uid, version, name, calories, proteins, fats, carbs
) VALUES (
	$1, (SELECT MAX(version) + 1 FROM products WHERE uid = $1),
	$2, $3, $4, $5, $6
) RETURNING version, created_at, updated_at;`

const createProductQuery /* sql */ = `INSERT INTO products (
	uid, version, name, calories, proteins, fats, carbs
) VALUES (
	$1, 1, $2, $3, $4, $5, $6
) RETURNING version, created_at, updated_at;`

const deleteByRefQuery /* sql */ = `DELETE FROM products WHERE uid = $1 AND version = $2;`

const updateProductQuery /* sql */ = `UPDATE products
SET name = $3, calories = $4, proteins = $5,
	fats = $6, carbs = $7
WHERE uid = $1 AND version = $2
RETURNING uid, version, created_at, updated_at;`

const searchProductsQuery /* sql */ = "SELECT " + productColumns + ` FROM products
WHERE name ILIKE '%' || $1 || '%'`

const nameTakenQuery /* sql */ = "SELECT 1 FROM products WHERE name LIKE $1"

const maxProductVersionQuery /* sql */ = "SELECT COALESCE(MAX(version), 0) FROM products WHERE uid = $1"

const usedInDishesQuery /* sql */ = `SELECT 1
FROM ingredients
WHERE ingredientable_type = $1
	AND ingredientable_uid = $2
	AND ingredientable_version = $3
LIMIT 1;`
