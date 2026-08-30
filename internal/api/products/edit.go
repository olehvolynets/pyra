package products

import (
	"database/sql"
	"errors"
	"net/http"

	"pyra/internal/api/base"
	"pyra/pkg/nutrition"
)

type EditProductHandler struct {
	*base.Handler
	productRepo nutrition.ProductRepository
}

func (h *EditProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.RequestLogger(r)

	ref, err := productRef(r)
	if err != nil {
		log.DebugContext(ctx, "malformed product UID or version", "error", err, "path", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.productRepo.FindByRef(ctx, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.DebugContext(ctx, "product not found", "uid", ref.UID, "version", ref.Version)
			h.NotFound(w, r)
			return
		}

		log.ErrorContext(ctx, "failed to retrieve a record", "error", err)
		h.InternalServerError(w)
		return
	}

	form := FormFromProduct(product)

	h.Render(w, r, "edit-product", form)
}
