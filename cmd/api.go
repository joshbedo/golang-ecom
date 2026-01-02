package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// Tips for writing better applications: https://12factor.net/

type application struct {
	config config
	// logger
	// db driver
	// email driver
}

// Dev/Prod should have close parity
type config struct {
	addr string // :8000
	db   dbConfig
}

type dbConfig struct {
	dsn string // postgres://user:password@host:port/dbname
}

// Mount -> routes
// Agnostic could use gorilla, chi, gin, etc. since it uses http.Handler
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// user -> handler GET /products -> service getProducts -> repo SELECT * FROM products -> throw error

	// Middleware
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // important for rate limiting and analytics or tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	// Set a timeout value on request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all gravy"))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi."))
	})

	// http.ListenAndServe(":3333", r)

	return r
}

// Run server
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: h,
		// @todo: Could be moved to config
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
