package products

import (
	"database/sql"
	"errors"
	"net/http"

	"pyra/internal/api/handler"
)

func DeleteProduct(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := handler.RequestLogger(r)

	ref, err := productRef(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// FIX: no need, handle NotFound in deletion
	_, err = api.ProductRepo.FindByRef(ctx, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		log.ErrorContext(ctx, "failed to find product", "error", err)
		handler.InternalServerError(w)
		return
	}

	err = api.ProductRepo.Delete(ctx, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		log.ErrorContext(ctx, "failed to delete product", "error", err)
		handler.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
