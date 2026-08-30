package nutrition

import (
	"context"

	"pyra/pkg/db"
)

type DishRef struct {
	UID DishUID
	Version DishVersion
}

type DishRepository interface {
	db.Repository[DishRepository]

	Index(context.Context) ([]Dish, error)

	FindByRef(context.Context, DishRef) (Dish, error)
	FindAllByRefs(context.Context, []DishRef) ([]Dish, error)

	Versions(context.Context, DishUID) ([]Dish, error)
	FindAllByProductRef(context.Context, ProductRef) ([]Dish, error)

	IsNameTaken(context.Context, DishName, DishUID) (bool, error)

	Create(context.Context, *Dish) error
	Delete(context.Context, DishRef) error
}
