package auth

import (
	"net/http"

	"pyra/internal/api/handler"
)

func SignOut(api *API, w http.ResponseWriter, r *http.Request) {
	session := handler.RequestSession(r)

	delete(session.Values, handler.UserIDSessionKey)

	w.Header().Add("HX-Redirect", "/signIn")
	w.WriteHeader(http.StatusNoContent)
}
