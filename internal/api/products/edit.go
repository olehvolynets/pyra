package products

import (
	"database/sql"
	"errors"
	"net/http"
	"pyra/internal/api/handler"
)

func EditProduct(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := handler.RequestLogger(r)

	ref, err := productRef(r)
	if err != nil {
		log.DebugContext(ctx, "malformed product UID or version", "error", err, "path", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := api.ProductRepo.FindByRef(ctx, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.DebugContext(ctx, "product not found", "uid", ref.UID, "version", ref.Version)
			handler.NotFound(w, editProductTemplate)
			return
		}

		log.ErrorContext(ctx, "failed to retrieve a record", "error", err)
		handler.InternalServerError(w)
		return
	}

	form := FormFromProduct(product)

	handler.Render(w, editProductTemplate, "edit-product", form)
}
