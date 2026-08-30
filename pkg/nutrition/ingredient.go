package nutrition

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidUnit    = errors.New("invalid value for a measurement unit")
	ErrNegativeAmount = errors.New("ingredient amount can't be less than 0")
)

type IngredientType uint16

const (
	IngredientNone IngredientType = iota
	IngredientProduct
	IngredientDish
)

type Ingredient struct {
	IngredientRecord

	Errors IngredientErrors
}

func (i *Ingredient) Macro() Macro {
	return i.Ingredientable.Macro.Scale(i.Amount.Float())
}

type IngredientRecord struct {
	DishUID     DishUID     `fake:"-"`
	DishVersion DishVersion `fake:"-"`

	Ingredientable Ingredientable

	Amount Measurement
	Unit   MeasurementUnit

	Idx uint16 // Defines presentation order for forms, recepies, etc.
}

// Ingredientable - item used as an ingredient (can be either product, or another dish, etc.).
type Ingredientable struct {
	Type    IngredientType ``
	UID     UID            ``
	Version Version        ``

	Name string

	Macro
}

type IngredientErrors struct {
	Dish error

	Ingredientable     error
	IngredientableType error

	Amount error
	Unit   error

	Idx error
}

// func NewIngredient(
// 	dishID DishID,
// 	ingredientableID uint64,
// 	ingredientableType IngredientableType,
// 	amt float32,
// 	unit MeasurementUnit,
// ) (Ingredient, error) {
// 	if unit == InvalidUnit {
// 		return Ingredient{}, fmt.Errorf("%w: %v", ErrInvalidUnit, unit)
// 	}
//
// 	if amt < 0 {
// 		return Ingredient{}, ErrNegativeAmount
// 	}
//
// 	return Ingredient{
// 		DishID:           dishID,
// 		IngredientableID: ingredientableID,
// 		Amount:           amt,
// 		Unit:             unit,
// 	}, nil
// }

func (it IngredientType) String() string {
	switch it {
	case IngredientNone:
		return "none"
	case IngredientProduct:
		return "product"
	case IngredientDish:
		return "dish"
	default:
		panic("WTF?")
	}
}

const ingredientableDebugFormat = `
Ref: %s/%d
Type: %s
Calories: %.2f
Proteins: %.2f
Fats: %.2f
Carbs: %.2f
`

func (ing Ingredientable) String() string {
	return fmt.Sprintf(ingredientableDebugFormat, ing.UID, ing.Version, ing.Type,
		ing.Calories.Float(), ing.Proteins.Float(), ing.Fats.Float(), ing.Carbs.Float())
}
