package middleware

import "net/http"

func Authenticated(next http.Handler) http.Handler {
	return next
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	l := log.FromContext(r.Context())
	// 	s := session.FromContext(r.Context())
	//
	// 	if _, ok := s.Values[UserIDSessionKey]; ok {
	// 		next.ServeHTTP(w, r)
	// 	} else {
	// 		l.Trace("unauthenticated access to a protected endpoint")
	// 		// s.AddFlash("Please sign in first")
	//
	// 		http.Redirect(w, r, "/signIn", http.StatusFound)
	// 	}
	// })
}
