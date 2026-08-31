package handler

import (
	"errors"
	"fmt"
	"html/template"

	"net/http"

	"pyra/pkg/log"
	"pyra/pkg/session"
)

var ErrNoUsesr = errors.New("no current user")

type HandlerFn[T any] func(c *T, w http.ResponseWriter, r *http.Request)

func New[T any](ctx *T, f HandlerFn[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(ctx, w, r)
	}
}

func RequestSession(r *http.Request) *session.Session {
	return session.FromContext(r.Context())
}

func RequestLogger(r *http.Request) *log.Logger {
	return log.FromContext(r.Context())
}

func Render(w http.ResponseWriter, t *template.Template, name string, data any) {
	if err := t.ExecuteTemplate(w, name, data); err != nil {
		panic(fmt.Errorf("render failed: %w", err))
	}
}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func NotFound(w http.ResponseWriter, t *template.Template) {
	w.WriteHeader(http.StatusNotFound)
	Render(w, t, "not-found-error", nil)
}

func IsAuthenticated(r *http.Request) bool {
	return true
	log := log.FromContext(r.Context())
	s := session.FromContext(r.Context())

	userId, ok := s.Values[UserIDSessionKey]

	if ok {
		log.Debug("USER_ID", "id", userId.(uint64))
	}

	return ok
}

// func CurrentUser(r *http.Request) (*auth.User, error) {
// 	s := h.Session(r)
// 	userID, ok := s.Values[UserIDSessionKey]
// 	if !ok {
// 		return nil, ErrNoUsesr
// 	}
//
// 	user, err := h.userRepo.FindByID(r.Context(), userID.(uint64))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &user, nil
// }
