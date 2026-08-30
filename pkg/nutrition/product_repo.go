package nutrition

import (
	"context"
	"time"

	"pyra/pkg/db"
)

type ProductRepository interface {
	db.Repository[ProductRepository]

	Index(context.Context) ([]Product, error)
	FindAllByRefs(context.Context, []ProductRef) ([]Product, error)
	ForDish(context.Context, DishUID, DishVersion) ([]Product, error)

	FindByRef(context.Context, ProductRef) (Product, error)

	Versions(context.Context, ProductUID) ([]Product, error)

	Create(context.Context, *Product) error
	CreateVersion(context.Context, *Product) error
	Delete(context.Context, ProductRef) error
	Update(context.Context, *Product) error
	Archive(context.Context, ProductRef, time.Time) error

	CountAll(context.Context) (int, error)

	IsNameTaken(context.Context, ProductName) (bool, error)
	UsedInDishes(context.Context, ProductRef) (bool, error)
	MaxVersion(context.Context, ProductUID) (ProductVersion, error)

	Search(ctx context.Context, searchStr string) ([]Product, error)
}
