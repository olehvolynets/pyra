package dishes

import (
	"database/sql"
	"errors"
	"net/http"

	"pyra/internal/api/base"
	"pyra/pkg/nutrition"
)

type ShowDishHandler struct {
	*base.Handler

	DishRepo nutrition.DishRepository
	ProductRepo nutrition.ProductRepository
}

func (h *ShowDishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.RequestLogger(r)

	ref, err := dishRef(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dish, err := h.DishRepo.FindByRef(ctx, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.NotFound(w, r)
		} else {
			log.TraceContext(ctx, "failed to find dish", "error", err)
			h.InternalServerError(w)
		}

		return
	}

	versions, err := h.DishRepo.Versions(ctx, dish.UID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.TraceContext(ctx, "failed to retrieve dish versions", "error", err)
		h.InternalServerError(w)
		return
	}

	products, err := h.ProductRepo.ForDish(ctx, dish.UID, dish.Version)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.TraceContext(ctx, "failed to retrieve dishe's products", "error", err)
		h.InternalServerError(w)
		return
	}

	h.Render(w, r, "dish-details", DishDetails{dish, versions, products})
}

type DishDetails struct {
	Dish     nutrition.Dish
	Versions []nutrition.Dish
	Products []nutrition.Product
}
