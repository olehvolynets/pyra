package ingredients

import (
	"context"
	"fmt"
	"pyra/pkg/db"
	"pyra/pkg/nutrition"
	"strings"
)

type Repository struct {
	db db.DBTX
}

func NewRepository(db db.DBTX) nutrition.IngredientRepository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) BeginTx(ctx context.Context) (db.DBTX, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *Repository) WithTx(tx db.DBTX) nutrition.IngredientRepository {
	return &Repository{
		db: tx,
	}
}

func (r *Repository) LoadIngredientables(
	ctx context.Context,
	ingredients []nutrition.Ingredient,
) error {
	if len(ingredients) == 0 {
		return nil
	}

	type mapping struct {
		Indexes []int
		Placeholders []string
	}

	partition := make(map[nutrition.IngredientType]mapping, 2)
	values := make([]any, 0, len(ingredients) * 2) // len * (UID + version)

	for i, ing := range ingredients {
		switch ing.Ingredientable.Type {
		case nutrition.IngredientProduct, nutrition.IngredientDish:
			part := partition[ing.Ingredientable.Type]

			part.Indexes = append(part.Indexes, i)

			placeholder := fmt.Sprintf("($%d, $%d)", i * 2 + 1, i * 2 + 2)
			part.Placeholders = append(part.Placeholders, placeholder)
			values = append(values, ing.Ingredientable.UID, ing.Ingredientable.Version)

			partition[ing.Ingredientable.Type] = part
		default:
			panic(fmt.Errorf("unhandled ingredientable type %d", ing.Ingredientable.Type))
		}
	}

	const productsQuery = `
SELECT uid, version, calories, proteins, fats, carbs, %s AS ingredientable_type
FROM products
WHERE (uid, version) IN (
	%s
)`
	const dishesQuery = `
SELECT uid, version, calories, proteins, fats, carbs, %s AS ingredientable_type
FROM dishes
WHERE (uid, version) IN (
	%s
)`
	const union = "\nUNION\n"

	var query strings.Builder

	if productPart, ok := partition[nutrition.IngredientProduct]; ok {
		fmt.Fprintf(&query, productsQuery,
			nutrition.IngredientProduct, strings.Join(productPart.Placeholders, ",\n\t"))
	}
	
	if dishPart, ok := partition[nutrition.IngredientDish]; ok {
		if query.Len() > 0 {
			query.WriteString(union)
		}

		fmt.Fprintf(&query, dishesQuery,
			nutrition.IngredientDish, strings.Join(dishPart.Placeholders, ",\n\t"))
	}

	fmt.Fprint(&query, ";")

	rows, err := r.db.QueryContext(ctx, query.String(), values...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var uid nutrition.UID
		var version nutrition.Version
		var t nutrition.IngredientType
		var m nutrition.Macro

		err = rows.Scan(&uid, &version, &m.Calories, &m.Proteins, &m.Fats, &m.Carbs, &t)
		if err != nil {
			return err
		}

		switch t {
		case nutrition.IngredientProduct, nutrition.IngredientDish:
			for _, i := range partition[t].Indexes {
				item := &ingredients[i]
				if item.Ingredientable.UID == uid && item.Ingredientable.Version == version {
					item.Ingredientable.Macro = m
				}
			}
		default:
			panic("WTF?")
		}
	}
	
	return nil
}

func (r *Repository) CreateIngredients(ctx context.Context, ingredients []nutrition.Ingredient) error {
	if len(ingredients) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(ingredients))
	values := make([]any, 0, len(ingredients)*8)

	for idx, ing := range ingredients {
		phIdx := idx * 8

		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			phIdx + 1, phIdx + 2, phIdx + 3, phIdx + 4, phIdx + 5, phIdx + 6, phIdx + 7, phIdx + 8,
		))

		values = append(values,
			ing.DishUID, ing.DishVersion,
			ing.Ingredientable.Type, ing.Ingredientable.UID, ing.Ingredientable.Version,
			ing.Amount, ing.Unit, ing.Idx)
	}

	query := `
INSERT INTO ingredients (
	dish_uid, dish_version, ingredientable_type, ingredientable_uid, ingredientable_version, amount, unit, idx
) VALUES %s`

	query = fmt.Sprintf(query, strings.Join(placeholders, ",\n\t"))

	result, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != int64(len(ingredients)) {
		return fmt.Errorf("expected to insert %d ingredients, got %d", len(ingredients), affected)
	}

	return nil
}
