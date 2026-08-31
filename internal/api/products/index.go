package products

import (
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/nutrition"
)

type ListProductsRenderContext struct {
	handler.RenderContext

	Products []nutrition.Product
	Form ProductForm
}

func ListProducts(api *API, w http.ResponseWriter, r *http.Request) {
	log := handler.RequestLogger(r)

	products, err := nutrition.ListProducts(r.Context(), api.ProductRepo)
	if err != nil {
		log.Error("failed to list produces", "error", err)
		handler.InternalServerError(w)
		return
	}

	details := ListProductsRenderContext{
		Flashes: []handler.FlashMessage{
			{ Severity: handler.FlashInfo, Msg: "Info message" },
			{ Severity: handler.FlashWarning, Msg: "Warning message" },
			{ Severity: handler.FlashError, Msg: "Error message" },
		},
		Products: products,
		Form: ProductForm{Per: "100"},
	}

	handler.Render(w, productListTemplate, "product-list", details)
}
