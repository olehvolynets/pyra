// Package api - root Pyra mux.
package api

import (
	"net/http"

	"pyra/internal/api/handler"
	"pyra/internal/api/middleware"
	"pyra/internal/api/dishes"
	"pyra/internal/api/products"
	"pyra/pkg/db"
	"pyra/pkg/log"
)

func Mux(db db.DBTX, l *log.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	// authAPI := auth.NewAPI(baseAPI)
	productsAPI := products.NewAPI(db)
	dishesAPI := dishes.NewAPI(db)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if handler.IsAuthenticated(r) {
				http.Redirect(w, r, "/products", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			}
		} else {
			http.NotFound(w, r)
		}
	})

	Authenticated := middleware.Authenticated

	// mux.Handle("GET /login", authAPI.SignIn())
	// mux.Handle("GET /auth/google", authAPI.GoogleAuthorize())
	// mux.Handle("GET /auth/google/callback", authAPI.GoogleCallback())
	// mux.Handle("GET /logout", Authenticated(authAPI.SignOut()))

	mux.Handle("GET /products", Authenticated(productsAPI.List()))
	mux.Handle("GET /products/{uid}/{version}", Authenticated(productsAPI.Show()))
	mux.Handle("GET /products/{uid}/{version}/edit", Authenticated(productsAPI.Edit()))
	mux.Handle("POST /products", Authenticated(productsAPI.Create()))
	mux.Handle("PUT /products/{uid}/{version}", Authenticated(productsAPI.Update()))
	mux.Handle("DELETE /products/{uid}/{version}", Authenticated(productsAPI.Delete()))
	mux.Handle("POST /products/search", Authenticated(productsAPI.Search()))

	mux.Handle("GET /dishes", Authenticated(dishesAPI.Index()))
	mux.Handle("GET /dishes/{id}", Authenticated(dishesAPI.Show()))
	mux.Handle("GET /dishes/new", Authenticated(dishesAPI.New()))
	mux.Handle("POST /dishes", Authenticated(dishesAPI.Create()))

	return mux
}
