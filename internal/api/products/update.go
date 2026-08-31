package products

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/nutrition"
)

func UpdateProduct(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := handler.RequestLogger(r)
	session := handler.RequestSession(r)

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
		handler.Render(w, editProductTemplate, "edit-product", form)
		return
	}

	product.ProductRef = ref

	err = nutrition.UpdateProduct(r.Context(), api.ProductRepo, &product)
	if err != nil {
		errMsg := fmt.Sprintf("couldn't update the product: %s", err.Error())
		log.DebugContext(ctx, errMsg)

		if errors.Is(err, sql.ErrNoRows) {
			log.TraceContext(ctx, "product not found", "error", err)
			handler.NotFound(w, editProductTemplate)
			return
		}

		// TODO: handle different kinds of errors since lost DB connection should result in 500.
		session.AddFlash(errMsg)
		w.WriteHeader(http.StatusUnprocessableEntity)
		handler.Render(w, editProductTemplate, "edit-product", form)
		return
	}

	loc := fmt.Sprintf("/products/%s/%d", product.UID, product.Version)
	w.Header().Add("HX-Redirect", loc)
	// http.Redirect(w, r, loc, http.StatusFound)
}
