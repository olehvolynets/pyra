package dishes

import (
	"io"
	"net/http"

	"pyra/internal/api/handler"
	// "pyra/pkg/nutrition"
)

func NewDish(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := handler.RequestLogger(r)

	if err := r.ParseForm(); err != nil {
		l.TraceContext(ctx, "failed to parse params", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Error parsing params")
		return
	}

	// productIDstr := r.PostForm["product"]
	// productAmounts := r.PostForm["amounts"]
	//
	// if len(productIDstr) != len(productAmounts) {
	// 	// TODO: log something
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	io.WriteString(w, "Number of product IDs doesn't match the number of amounts")
	// 	return
	// }
	//
	// ingredients := make([]nutrition.Ingredient, len(productIDstr))
	// productIDs := make([]nutrition.ProductID, len(productIDstr))
	//
	// for i := 0; i < len(productIDstr); i++ {
	// 	id, err := strconv.ParseUint(productIDstr[i], 10, 64)
	// 	if err != nil {
	// 	}
	//
	// 	amount, err := strconv.ParseFloat(productIDstr[i], 64)
	// 	if err != nil {
	// 	}
	//
	// 	ingredients[i], err = nutrition.NewIngredient(0, id, float32(amount), int32(nutrition.Gramms))
	// 	if err != nil {
	// 	}
	//
	// 	productIDs[i] = ingredients[i].ProductID
	// }
	//
	// products, err := h.productRepo.FindAllByIDs(ctx, productIDs)
	// if len(products) != len(productIDs) {
	// 	l.DebugContext(ctx, "some of the requested products weren't found")
	//
	// 	form := DishForm{}
	// 	w.WriteHeader(http.StatusNotFound)
	// 	h.Render(w, r, "new-dish", form)
	// 	return
	// }
	//
	// dish := nutrition.Dish{Name: r.PostForm.Get("name")}
	// dish.SetIngredients(ingredients)
	//
	// for i := 0; i < len(productIDs); i++ {
	// 	productId := productIDs[i]
	// 	amount := productAmounts[i]
	//
	// 	_ = productId
	// 	_ = amount
	// }
	//
	// id, err := h.svc.Create(r.Context(), dish)
	// if err != nil {
	// 	form := DishForm{}
	// 	w.WriteHeader(http.StatusUnprocessableEntity)
	// 	h.Render(w, r, "new-dish", form)
	// }
	//
	// http.Redirect(w, r, fmt.Sprintf("/dishes/%d", id), http.StatusFound)
}
