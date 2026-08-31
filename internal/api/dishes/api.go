package dishes

import (
	"net/http"

	"pyra/internal/api/base"
	"pyra/internal/dishes"
	"pyra/internal/products"
)

type API struct {
	*base.API
}

// NewAPI - creates API instance for dishes-related endpoints.
func NewAPI(api *base.API) *API {
	return &API{
		API: api,
	}
}

func (api *API) Index() http.Handler {
	baseHandler := api.NewHandler("view/dishes/index.html")

	return &ListDishesHandler{
		Handler: baseHandler,
		repo:    dishes.NewRepository(api.DB),
	}
}

func (api *API) Show() http.Handler {
	baseHandler := api.NewHandler("view/dishes/show.html")

	return &ShowDishHandler{
		Handler:     baseHandler,
		DishRepo:    dishes.NewRepository(api.DB),
		ProductRepo: products.NewRepository(api.DB),
	}
}

func (api *API) New() http.Handler {
	baseHandler := api.NewHandler("view/dishes/new.html")

	return &NewDishHandler{
		Handler: baseHandler,
	}
}

func (api *API) Create() http.Handler {
	baseHandler := api.NewHandler("view/dishes/new.html")

	return &CreateDishHandler{
		Handler: baseHandler,
		DishRepo: dishes.NewRepository(api.DB),
		// IngredientRepo: dishes.NewRepository(api.DB),
	}
}

func (api *API) Delete() http.Handler {
	return &DeleteDishHandler{
		Handler: api.NewHandler(),
		DishRepo: dishes.NewRepository(api.DB),
	}
}
