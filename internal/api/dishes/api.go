package dishes

import (
	"embed"
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

//go:embed templates
var templateFS embed.FS

func init() {
	handler.GlobalAddFuncs(URIHelpers)
	dishListTemplate = handler.ExtendedTemplate(templateFS, "templates/index.html")
	dishTemplate = handler.ExtendedTemplate(templateFS, "templates/show.html")
	newDishTemplate = handler.ExtendedTemplate(templateFS, "templates/new.html")
	// editDishTemplate = handler.ExtendedTemplate(templateFS, "templates/edit.html")
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
