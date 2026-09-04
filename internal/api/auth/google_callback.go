package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"pyra/internal/api/handler"
	"pyra/pkg/auth"
)

func GoogleCallback(api *API, w http.ResponseWriter, r *http.Request) {
	log := handler.RequestLogger(r)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		log.ErrorContext(r.Context(), "invalid code")
		handler.InternalServerError(w)
		return
	}

	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, httpClient)

	tok, err := googleConfig.Exchange(ctx, code)
	if err != nil {
		log.ErrorContext(ctx, "failed to exchange access grant", "error", err)
		handler.InternalServerError(w)
		return
	}

	client := googleConfig.Client(ctx, tok)

	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.ErrorContext(ctx, "failed to fetch user details from Google", "error", err)
		handler.InternalServerError(w)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.ErrorContext(
			ctx,
			fmt.Sprintf("Google responded with a %d trying to fetch user information", response.StatusCode),
		)
		handler.InternalServerError(w)
		return
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.ErrorContext(ctx, "failed to read response", "error", err)
		handler.InternalServerError(w)
		return
	}

	var u auth.GoogleUser
	if err := json.Unmarshal(responseBytes, &u); err != nil {
		log.ErrorContext(ctx, "failed to unmarshal response", "error", err)
		handler.InternalServerError(w)
		return
	}

	log.Inspect(u)

	// FIX: passing nil instead of db.DBTX. Should not pass it at all, only repositories.
	authSvc := auth.NewService(nil, api.ProviderRepo, api.UserRepo)
	user, err := authSvc.SignIn(r.Context(), u)
	if err != nil {
		log.ErrorContext(ctx, "sign in failed", "error", err)
		handler.InternalServerError(w)
		return
	}

	session := handler.RequestSession(r)
	session.Values[handler.UserIDSessionKey] = user.ID

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}
