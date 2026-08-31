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
	flashes := make([]handler.FlashMessage, 0)

	err := r.ParseForm()
	if err != nil {
		// FIX: use immediate flashes instead of session since there's no redirect.
		session.AddFlash(fmt.Sprintf("failed to parse form: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	form := NewProductForm(r.FormValue)
	product := form.BuildProduct()

	if form.HasErrors() {
		log.DebugContext(ctx, "create product validation error", "error", form.Errors)
	} else {
		create_err := nutrition.CreateProduct(ctx, api.ProductRepo, &product)
		if create_err != nil {
			log.DebugContext(ctx, "failed to save product", "error", create_err)

			if errors.Is(create_err, nutrition.ErrProductInvalid) {
				product.Errors.Base = create_err
			} else {
				handler.InternalServerError(w)
				return
			}
		}
	}

	if product.HasErrors() {
		form.SetProductErrors(product.Errors)

		if form.Errors.Base != nil {
			flashes = append(flashes, handler.NewFlash(handler.FlashError, fmt.Sprintf("failed to create a product: %s", form.Errors.Base)))
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		form = ProductForm{Per: "100"}
	}

	products, err := nutrition.ListProducts(r.Context(), api.ProductRepo)
	if err != nil {
		log.Error("failed to list produces", "error", err)
		handler.InternalServerError(w)
		return
	}

	details := ListProductsRenderContext{
		Flashes:  flashes,
		Products: products,
		Form:     form,
	}

	handler.Render(w, productListTemplate, "product-list", details)
}
