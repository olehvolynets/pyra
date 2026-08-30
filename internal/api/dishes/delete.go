package dishes

import (
	"database/sql"
	"errors"
	"net/http"

	"pyra/internal/api/base"
	"pyra/pkg/log"
	"pyra/pkg/nutrition"
)

type DeleteDishHandler struct {
	*base.Handler

	DishRepo nutrition.DishRepository
}

func (h *DeleteDishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.FromContext(ctx)

	ref, err := dishRef(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = nutrition.DeleteDish(ctx, h.DishRepo, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.DebugContext(ctx, "failed to delete dish", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
