package main

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/handlers"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	datastore := cache.Create()
	handler := handlers.New(datastore)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handler.Other)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{typeMetric}/{nameMetric}/{value}", handler.UpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{typeMetric}/{nameMetric}", handler.GetMetric)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
