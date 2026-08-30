package dishes

import (
	"fmt"
	"net/http"

	"pyra/internal/api/base"
	"pyra/pkg/log"
	"pyra/pkg/nutrition"
)

type CreateDishHandler struct {
	*base.Handler

	DishRepo       nutrition.DishRepository
	IngredientRepo nutrition.IngredientRepository
}

func (h *CreateDishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.FromContext(ctx)
	session := h.Session(r)

	if err := r.ParseForm(); err != nil {
		session.AddFlash(fmt.Sprintf("failed to parse form: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	form := NewDishForm(r.FormValue)
	dish := form.BuildDish()

	if form.HasErrors() {
		log.DebugContext(ctx, "create dish validation error", "error", form.Errors)
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(w, r, "new-dish", form)
		return
	}

	if err := nutrition.CreateDish(ctx, h.DishRepo, h.IngredientRepo, &dish); err != nil {
		log.DebugContext(ctx, "failed to save dish", "error", err)
		h.InternalServerError(w)
		return
	}

	http.Redirect(w, r, DishURI(dish), http.StatusFound)
}
