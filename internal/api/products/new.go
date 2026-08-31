package products

import (
	"net/http"

	"pyra/internal/api/handler"
)

func NewProduct(api *API, w http.ResponseWriter, r *http.Request) {
	form := ProductForm{Per: "100"}

	handler.Render(w, newProductTemplate, "new-product", form)
}
