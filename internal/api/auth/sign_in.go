package auth

import (
	"net/http"

	"pyra/internal/api/handler"
)

func SignIn(api *API, w http.ResponseWriter, r *http.Request) {
	handler.Render(w, signInTemplate, "sign-in", nil)
}
