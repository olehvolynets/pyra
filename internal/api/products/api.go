package products

import (
	"html/template"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/internal/dishes"
	"pyra/internal/products"
	"pyra/pkg/db"
	"pyra/pkg/nutrition"
)

type API struct {
	ProductRepo nutrition.ProductRepository
	DishRepo    nutrition.DishRepository
}

var (
	productListTemplate *template.Template
	productTemplate     *template.Template
	newProductTemplate  *template.Template
	editProductTemplate *template.Template
)

func init() {
	handler.GlobalAddFuncs(URIHelpers)
	productListTemplate = handler.ExtendedTemplate("view/products/index.html")
	productTemplate = handler.ExtendedTemplate("view/products/show.html")
	newProductTemplate = handler.ExtendedTemplate("view/products/new.html")
	editProductTemplate = handler.ExtendedTemplate("view/products/edit.html")
}

func NewAPI(db db.DBTX) *API {
	repo := products.NewRepository(db)
	dishRepo := dishes.NewRepository(db)

	return &API{
		ProductRepo: repo,
		DishRepo:    dishRepo,
	}
}

func (productApi *API) List() http.Handler {
	return handler.New(productApi, ListProducts)
}

func (productApi *API) Show() http.Handler {
	return handler.New(productApi, ShowProduct)
}

func (productApi *API) New() http.Handler {
	return handler.New(productApi, NewProduct)
}

func (productApi *API) Create() http.Handler {
	return handler.New(productApi, CreateProduct)
}

func (productApi *API) Edit() http.Handler {
	return handler.New(productApi, EditProduct)
}

func (productApi *API) Update() http.Handler {
	return handler.New(productApi, UpdateProduct)
}

func (productApi *API) Delete() http.Handler {
	return handler.New(productApi, DeleteProduct)
}

func (productApi *API) Search() http.Handler {
	return handler.New(productApi, SearchProduct)
}
