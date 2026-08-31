package dishes

import (
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/nutrition"
)

func ListDishes(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := handler.RequestLogger(r)

	dishes, err := nutrition.ListDishes(ctx, api.DishRepo)
	if err != nil {
		log.ErrorContext(ctx, "failed to list dishes", "error", err)
		handler.InternalServerError(w)
		return
	}

	handler.Render(w, dishListTemplate, "dish-list", dishes)
}
