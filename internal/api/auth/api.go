package auth

import (
	"html/template"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/auth"
)

var signInTemplate *template.Template

func init() {
	signInTemplate = handler.ExtendedTemplate("view/auth/sign_in.html")
}

type API struct {
	UserRepo auth.UserRepository
	ProviderRepo auth.ProviderRepository
}

func (authAPI *API) SignIn() http.Handler {
	return handler.New(authAPI, SignIn)
}

func (authAPI *API) SignOut() http.Handler {
	return handler.New(authAPI, SignOut)
}

func (authAPI *API) GoogleAuthorize() http.Handler {
	return handler.New(authAPI, GoogleAuthorize)
}

func (authAPI *API) GoogleCallback() http.Handler {
	return handler.New(authAPI, GoogleCallback)
}
