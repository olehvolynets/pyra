package server

import (
	"net/http"
	"slices"
	"strings"
	"time"
	"uuid"

	"pyra/pkg/log"
)

type responseStatusCatcher struct {
	http.ResponseWriter
	statusCode int
}

func (c *responseStatusCatcher) WriteHeader(statusCode int) {
	c.statusCode = statusCode
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *responseStatusCatcher) StatusCode() int {
	if c.statusCode == 0 {
		return http.StatusOK
	}

	return c.statusCode
}

func Logger(logger *log.Logger, f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startT := time.Now()

		traceID := uuid.New()
		l := logger.With("traceId", traceID.String())

		ctx := r.Context()
		ctxWithLog := log.CtxWithLogger(ctx, l)
		r = r.WithContext(ctxWithLog)

		ww := responseStatusCatcher{ResponseWriter: w}
		f.ServeHTTP(&ww, r)

		endT := time.Now()
		took := endT.Sub(startT)

		var level log.Level
		if ww.StatusCode() < 500 {
			level = log.LevelInfo
		} else {
			level = log.LevelError
		}

		filteredPaths := []string{
			"/apple-touch-icon.png",
			"/apple-touch-icon-precomposed.png",
			"/favicon.ico",
		}
		if slices.Contains(filteredPaths, r.URL.Path) || strings.HasPrefix(r.URL.Path, "/assets/") {
			return
		}

		l.Log(ctx, level, "",
			"event", log.RequestEvent,
			"status", ww.StatusCode(),
			"method", r.Method,
			"path", r.URL.Path,
			"took", took,
			"location", ww.Header().Get("Location"),
		)
	})
}
