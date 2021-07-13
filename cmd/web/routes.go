package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	config2 "github.com/liomazza/bookings/internal/config"
	handlers2 "github.com/liomazza/bookings/internal/handlers"
	"net/http"
)

func routes(app *config2.AppConfig) http.Handler {
	/*mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))*/

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers2.Repo.Home)
	mux.Get("/about", handlers2.Repo.About)
	mux.Get("/generals-quarters", handlers2.Repo.Generals)
	mux.Get("/majors-suite", handlers2.Repo.Majors)

	mux.Get("/search-availability", handlers2.Repo.Availability)
	mux.Post("/search-availability", handlers2.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers2.Repo.AvailabilityJSON)

	mux.Get("/contact", handlers2.Repo.Contact)

	mux.Get("/make-reservation", handlers2.Repo.Reservation)
	mux.Post("/make-reservation", handlers2.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers2.Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
