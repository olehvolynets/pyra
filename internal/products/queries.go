package products

// lang: SQL
const productColumns = "uid, version, name, calories, proteins, fats, carbs, created_at, updated_at"

// lang: SQL
const indexProductsQuery = `SELECT 
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

// lang: SQL
const findByRefQuery = "SELECT " + productColumns + ` FROM products
WHERE uid = $1 AND version = $2
LIMIT 1;`

// lang: SQL
const productVersionsQuery = "SELECT " + productColumns + ` FROM products
	WHERE uid = $1
	ORDER BY version DESC
	LIMIT 20;`

// lang: SQL
const findAllByIDsQuery = "SELECT " + productColumns + " FROM products WHERE id in $0;"

// lang: SQL
const productsForDishQuery = "SELECT " + productColumns + ` FROM ingredients
JOIN products
ON ingredients.ingredientable_uid = products.uid
	AND ingredients.ingredientable_version = products.version
	AND ingredients.ingredientable_type = 'product'
WHERE ingredients.dish_uid = $1 AND ingredients.dish_version = $2;`

// lang: SQL
const createProductVersionQuery = `
INSERT INTO products (
	uid, version, name, calories, proteins, fats, carbs
) VALUES (
	$1, (SELECT MAX(version) + 1 FROM products WHERE uid = $1),
	$2, $3, $4, $5, $6
) RETURNING version, created_at, updated_at;`

// lang: SQL
const createProductQuery = `INSERT INTO products (
	uid, version, name, calories, proteins, fats, carbs
) VALUES (
	$1, 1, $2, $3, $4, $5, $6
) RETURNING version, created_at, updated_at;`

// lang: SQL
const deleteByRefQuery = "DELETE FROM products WHERE uid = $1 AND version = $2;"

// lang: SQL
const updateProductQuery = `UPDATE products
SET name = $3, calories = $4, proteins = $5,
	fats = $6, carbs = $7
WHERE uid = $1 AND version = $2
RETURNING uid, version, created_at, updated_at;`

// lang: SQL
const searchProductsQuery = "SELECT " + productColumns + ` FROM products
WHERE name ILIKE '%' || $1 || '%'`

// lang: SQL
const nameTakenQuery = "SELECT 1 FROM products WHERE name LIKE $1"

// lang: SQL
const maxProductVersionQuery = "SELECT COALESCE(MAX(version), 0) FROM products WHERE uid = $1"

// lang: SQL
const usedInDishesQuery = `SELECT 1
FROM ingredients
WHERE ingredientable_type = $1
	AND ingredientable_uid = $2
	AND ingredientable_version = $3
LIMIT 1;`
