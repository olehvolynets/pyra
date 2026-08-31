package dishes

import (
	"html/template"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/internal/dishes"
	"pyra/internal/ingredients"
	"pyra/internal/products"
	"pyra/pkg/db"
	"pyra/pkg/nutrition"
)

type API struct {
	DishRepo       	nutrition.DishRepository
	ProductRepo 	nutrition.ProductRepository
	IngredientRepo 	nutrition.IngredientRepository
}

var (
	dishListTemplate *template.Template
	dishTemplate     *template.Template
	newDishTemplate  *template.Template
	editDishTemplate *template.Template
)

func init() {
	handler.GlobalAddFuncs(URIHelpers)
	dishListTemplate = handler.ExtendedTemplate("view/dishes/index.html")
	dishTemplate = handler.ExtendedTemplate("view/dishes/show.html")
	newDishTemplate = handler.ExtendedTemplate("view/dishes/new.html")
	// editDishTemplate = handler.ExtendedTemplate("view/dishes/edit.html")
}

func NewAPI(db db.DBTX) *API {
	return &API{
		DishRepo:    dishes.NewRepository(db),
		ProductRepo: products.NewRepository(db),
		IngredientRepo: ingredients.NewRepository(db),
	}
}

func (api *API) Index() http.Handler {
	return handler.New(api, ListDishes)
}

func (api *API) Show() http.Handler {
	return handler.New(api, ShowDish)
}

func (api *API) New() http.Handler {
	return handler.New(api, NewDish)
}

func (api *API) Create() http.Handler {
	return handler.New(api, CreateDish)
}

func (api *API) Delete() http.Handler {
	return handler.New(api, DeleteDish)
}
