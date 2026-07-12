package logger

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
)

func Middleware(next http.Handler) http.Handler {
	h := hlog.NewHandler(log)

	accessHandler := hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			hlog.
				FromRequest(r).
				Info().
				Stringer("url", r.URL).
				Int("status_code", status).
				Int("bytes", size).
				Dur("elapsed_ms", duration).
				Str("method", r.Method).
				Msgf("Request")
		},
	)

	return h(accessHandler(next))
}
