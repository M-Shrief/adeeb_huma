package router

import (
	"adeeb_huma/logger"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

var R *chi.Mux

func NewRouter() *chi.Mux {
	R = chi.NewRouter()
	return R
}

// Use standard middlewares.
func UseMiddlewares() {
	// Logging requests details
	R.Use(logger.Middleware)
	// Specifying content type
	R.Use(middleware.AllowContentType("application/json"))
	// CORS config
	R.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// Adding rate limit by IP
	R.Use(httprate.LimitByIP(1000, time.Minute))
	// Recovering from Panic and returns 500 server error,
	// and it restarts the server.
	R.Use(middleware.Recoverer)
}
