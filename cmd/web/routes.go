package main

import (
	"net/http"

	"github.com/Soyaib10/comfort-cocoon/internal/config"
	"github.com/Soyaib10/comfort-cocoon/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter() // a new mux 'chi', mux is a http handler broadly called multiplexer.

	mux.Use(middleware.Recoverer) // if some panic happens in handling requestes then this manages that.
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJson)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservaition-summary", handlers.Repo.ReservationSummary)

	mux.Get("/contact", handlers.Repo.Contact)



	fileServer := http.FileServer(http.Dir(http.Dir("./static/")))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}