package dishes

import (
	"database/sql"
	"errors"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/nutrition"
)

func DeleteDish(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := handler.RequestLogger(r)

	ref, err := dishRef(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = nutrition.DeleteDish(ctx, api.DishRepo, ref)
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
