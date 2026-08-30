package dishes

import (
	"pyra/pkg/nutrition"
	"strings"
)

type DishForm struct {
	Name        string

	Ingredients []IngredientForm

	Errors nutrition.DishErrors
}

func NewDishForm(fetch func(key string) string) (form DishForm) {
	fetchValue := func(key string) string {
		return strings.TrimSpace(fetch(key))
	}

	form.Name = fetchValue("name")

	// TODO: fetch ingredients

	return
}

func (f *DishForm) BuildDish() (dish nutrition.Dish) {
	dish.Name = nutrition.DishName(f.Name)

	// TODO: handle ingredients.
	// dish.Ingredients = make([]nutrition.Ingredient, len(f.Ingredients))
	// for _, formIng := range f.Ingredients {
	// 	dish.Ingredients = append(dish.Ingredients, formIng.BuildIngredient())
	// }

	return
}

func (f *DishForm) HasErrors() bool {
	return f.Errors.HasErrors()
}

type IngredientForm struct {
	Type    string
	UID     string
	Version string
	Amount  string
}

func (f *IngredientForm) BuildIngredient() (ingredient nutrition.Ingredient) {
	ingredient.Ingredientable.UID = nutrition.UID(f.UID)

	return
}
