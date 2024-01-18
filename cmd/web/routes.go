package main

import (
	"net/http"

	"github.com/Soyaib10/comfort-cocoon/pkg/config"
	"github.com/Soyaib10/comfort-cocoon/pkg/handlers"
	"github.com/bmizerany/pat"
)

func routes(app *config.AppConfig) http.Handler {
	mux := pat.New() // a new mux, mux is a http handler broadly called multiplexer.

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))


	return mux
}