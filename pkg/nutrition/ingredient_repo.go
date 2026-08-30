package nutrition

import (
	"context"
	"pyra/pkg/db"
)

type IngredientRepository interface {
	db.Repository[IngredientRepository]

	CreateIngredients(context.Context, []Ingredient) error
	// LoadIngredientables - loads data of the items used as ingredients and
	// populates the corresponding attributes of the provided Ingredient's in place.
	LoadIngredientables(context.Context, []Ingredient) error
}
