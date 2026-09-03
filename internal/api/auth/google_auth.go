package auth

import (
	"net/http"
	"os"
	"uuid"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"pyra/internal/api"
)

var googleConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_OAUTH2_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"openid",
	},
	Endpoint: google.Endpoint,
}

type GoogleAuthHandler struct {
	*api.Handler
}

func (h *GoogleAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := h.RequestLogger(r)
	session := h.Session(r)

	state := uuid.New()
	session.Values["state"] = state.String()
	url := googleConfig.AuthCodeURL(state.String())

	log.Debug("redirecting to Google for sign in")
	http.Redirect(w, r, url, http.StatusSeeOther)
}
