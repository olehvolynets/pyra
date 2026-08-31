package dishes

import (
	"database/sql"
	"errors"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/nutrition"
)

type DishDetails struct {
	Dish     nutrition.Dish
	Versions []nutrition.Dish
	Products []nutrition.Product
}

func ShowDish(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := handler.RequestLogger(r)

	ref, err := dishRef(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dish, err := api.DishRepo.FindByRef(ctx, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handler.NotFound(w, dishTemplate)
		} else {
			log.TraceContext(ctx, "failed to find dish", "error", err)
			handler.InternalServerError(w)
		}

		return
	}

	versions, err := api.DishRepo.Versions(ctx, dish.UID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.TraceContext(ctx, "failed to retrieve dish versions", "error", err)
		handler.InternalServerError(w)
		return
	}

	products, err := api.ProductRepo.ForDish(ctx, dish.UID, dish.Version)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.TraceContext(ctx, "failed to retrieve dishe's products", "error", err)
		handler.InternalServerError(w)
		return
	}

	handler.Render(w, dishTemplate, "dish-details", DishDetails{dish, versions, products})
}
