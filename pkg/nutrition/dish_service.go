package nutrition

import (
	"context"
	"errors"

	"pyra/pkg/db"

	"github.com/google/uuid"
)

var (
	ErrDishInvalid   = errors.New("Dish is invalid")
	ErrDishNameTaken = errors.New("Dish with such name already exists")
)

func ListDishes(ctx context.Context, repo DishRepository) ([]Dish, error) {
	return repo.Index(ctx)
}

func CreateDish(
	ctx context.Context,
	repo DishRepository,
	ingredientRepo IngredientRepository,
	dish *Dish,
) (err error) {
	dish.UID = DishUID(uuid.New().String())
	dish.Version = 1

	for i := range dish.Ingredients {
		dish.Ingredients[i].Idx = uint16(i + 1)
	}

	tx, err := repo.BeginTx(ctx)
	if err != nil {
		return
	}
	defer db.RollbackGuard(ctx, tx, &err)

	repo = repo.WithTx(tx)
	// ingredientRepo = ingredientRepo.WithTx(tx)

	nameIsTaken, err := repo.IsNameTaken(ctx, dish.Name, dish.UID)
	if err != nil {
		return
	} else if nameIsTaken {
		dish.Errors.Name = ErrDishNameTaken
		return ErrDishInvalid
	}

	// err = ingredientRepo.LoadIngredientables(ctx, dish.Ingredients)
	// if err != nil {
	// 	return
	// }

	// TODO: do this match in IngredientRepo.LoadIngredientables
	// if len(ingredientables) != len(info.Ingredients) {
	// 	return dish, errors, fmt.Errorf("requested %d ingredients, found only %d", len(info.Ingredients), len(ingredientables))
	// }

	for i, ing := range dish.Ingredients {
		dish.Ingredients[i].DishUID = dish.UID
		dish.Ingredients[i].DishVersion = dish.Version
		dish.Macro = dish.Macro.Add(ing.Macro())
	}

	err = repo.Create(ctx, dish)
	if err != nil {
		return
	}

	// err = ingredientRepo.CreateIngredients(ctx, dish.Ingredients)
	// if err != nil {
	// 	return
	// }

	return tx.Commit()
}

func DeleteDish(ctx context.Context, repo DishRepository, ref DishRef) error {
	return repo.Delete(ctx, ref)
}
