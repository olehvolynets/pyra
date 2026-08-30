package products

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"pyra/internal/api/base"
	"pyra/pkg/nutrition"
)

type UpdateProductHandler struct {
	*base.Handler
	ProductRepo nutrition.ProductRepository
	DishRepo    nutrition.DishRepository
}

func (h *UpdateProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.RequestLogger(r)
	session := h.Session(r)

	ref, err := productRef(r)
	if err != nil {
		log.DebugContext(ctx, "malformed product UID or version", "error", err, "path", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	form := NewProductForm(r.FormValue)
	product := form.BuildProduct()

	if form.HasErrors() {
		log.DebugContext(ctx, "update product form error", "error", form.Errors.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(w, r, "edit-product", form)
		return
	}

	product.ProductRef = ref

	err = nutrition.UpdateProduct(r.Context(), h.ProductRepo, &product)
	if err != nil {
		errMsg := fmt.Sprintf("couldn't update the product: %s", err.Error())
		log.DebugContext(ctx, errMsg)

		if errors.Is(err, sql.ErrNoRows) {
			log.TraceContext(ctx, "product not found", "error", err)
			h.NotFound(w, r)
			return
		}

		// TODO: handle different kinds of errors since lost DB connection should result in 500.
		session.AddFlash(errMsg)
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(w, r, "edit-product", form)
		return
	}

	loc := fmt.Sprintf("/products/%s/%d", product.UID, product.Version)
	w.Header().Add("HX-Redirect", loc)
	// http.Redirect(w, r, loc, http.StatusFound)
}
