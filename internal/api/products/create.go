package products

import (
	"errors"
	"fmt"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/nutrition"
)

func CreateProduct(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := handler.RequestLogger(r)
	session := handler.RequestSession(r)

	err := r.ParseForm()
	if err != nil {
		session.AddFlash(fmt.Sprintf("failed to parse form: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	form := NewProductForm(r.FormValue)
	product := form.BuildProduct()

	if form.HasErrors() {
		log.DebugContext(ctx, "create product validation error", "error", form.Errors)
		w.WriteHeader(http.StatusUnprocessableEntity)
		handler.Render(w, newProductTemplate, "new-product", form)
		return
	}

	err = nutrition.CreateProduct(ctx, api.ProductRepo, &product)
	if err != nil && !errors.Is(err, nutrition.ErrProductInvalid) {
		log.DebugContext(ctx, "failed to save product", "error", err)
		handler.InternalServerError(w)
		return
	}

	if product.HasErrors() {
		form.SetProductErrors(product.Errors)
		session.AddFlash("failed to create a product")
		w.WriteHeader(http.StatusUnprocessableEntity)
		handler.Render(w, newProductTemplate, "new-product", form)
		return
	}

	http.Redirect(w, r,
		fmt.Sprintf("/products/%s/%d", product.UID, product.Version),
		http.StatusFound)
}
