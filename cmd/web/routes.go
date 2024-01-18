package main

import (
	"net/http"

	"github.com/Soyaib10/comfort-cocoon/pkg/config"
	"github.com/Soyaib10/comfort-cocoon/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter() // a new mux 'chi', mux is a http handler broadly called multiplexer.
	mux.Use(middleware.Recoverer) // if some panic happens in handling requestes then this manages that.
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	return mux
}