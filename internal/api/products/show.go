package products

import (
	"database/sql"
	"errors"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/nutrition"
)

type ProductDetails struct {
	Product      nutrition.Product
	Versions     []nutrition.Product
	UsedInDishes []nutrition.Dish
}

// GET /products/:uid/:version
func ShowProduct(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := handler.RequestLogger(r)

	ref, err := productRef(r)
	if err != nil {
		log.ErrorContext(ctx, "malformed product UID or version", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := api.ProductRepo.FindByRef(ctx, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		log.ErrorContext(ctx, "failed to retrieve a record", "error", err)
		handler.InternalServerError(w)
		return
	}

	versions, err := api.ProductRepo.Versions(ctx, product.UID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.ErrorContext(ctx, "failed to retrieve a record", "error", err)
		handler.InternalServerError(w)
		return
	}

	usedInDishes, err := api.DishRepo.FindAllByProductRef(ctx, product.ProductRef)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.ErrorContext(ctx, "failed to retrieve a record", "error", err)
		handler.InternalServerError(w)
		return
	}

	details := ProductDetails{
		Product:      product,
		Versions:     versions,
		UsedInDishes: usedInDishes,
	}

	handler.Render(w, productTemplate, "product-details", details)
}
