package dishes

import (
	"fmt"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/log"
	"pyra/pkg/nutrition"
)

func CreateDish(api *API, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.FromContext(ctx)
	session := handler.RequestSession(r)

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
		handler.Render(w, newDishTemplate, "new-dish", form)
		return
	}

	if err := nutrition.CreateDish(ctx, api.DishRepo, api.IngredientRepo, &dish); err != nil {
		log.DebugContext(ctx, "failed to save dish", "error", err)
		handler.InternalServerError(w)
		return
	}

	http.Redirect(w, r, DishURI(dish), http.StatusFound)
}
