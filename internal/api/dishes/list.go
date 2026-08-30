package dishes

import (
	"net/http"

	"pyra/internal/api/base"
	"pyra/pkg/nutrition"
)

type ListDishesHandler struct {
	*base.Handler
	repo nutrition.DishRepository
}

func (h *ListDishesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.RequestLogger(r)

	dishes, err := nutrition.ListDishes(ctx, h.repo)
	if err != nil {
		log.ErrorContext(ctx, "failed to list dishes", "error", err)
		h.InternalServerError(w)
		return
	}

	h.Render(w, r, "dish-list", dishes)
}
