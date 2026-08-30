package products

import (
	"database/sql"
	"errors"
	"net/http"

	"pyra/internal/api/base"
	"pyra/pkg/nutrition"
)

type ProductHandler struct {
	*base.Handler

	productRepo nutrition.ProductRepository
	dishRepo    nutrition.DishRepository
}

// GET /products/:uid/:version
func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.RequestLogger(r)

	ref, err := productRef(r)
	if err != nil {
		log.ErrorContext(ctx, "malformed product UID or version", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.productRepo.FindByRef(ctx, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		log.ErrorContext(ctx, "failed to retrieve a record", "error", err)
		h.InternalServerError(w)
		return
	}

	versions, err := h.productRepo.Versions(ctx, product.UID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.ErrorContext(ctx, "failed to retrieve a record", "error", err)
		h.InternalServerError(w)
		return
	}

	usedInDishes, err := h.dishRepo.FindAllByProductRef(ctx, product.ProductRef)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.ErrorContext(ctx, "failed to retrieve a record", "error", err)
		h.InternalServerError(w)
		return
	}

	details := ProductDetails{
		Product:      product,
		Versions:     versions,
		UsedInDishes: usedInDishes,
	}

	h.Render(w, r, "product-details", details)
}

type ProductDetails struct {
	Product      nutrition.Product
	Versions     []nutrition.Product
	UsedInDishes []nutrition.Dish
}
