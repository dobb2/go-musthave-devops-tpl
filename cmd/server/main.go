package main

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/backup"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/handlers"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := config.CreateServerConfig()
	log.Println(cfg)
	r := chi.NewRouter()
	datastore := cache.Create()

	if cfg.StoreFile != "" {
		backup := backup.New(datastore)
		if cfg.Restore {
			backup.Restore(cfg)
		}

		c := make(chan struct{})

		if cfg.StoreInterval == 0 {
			datastore.AddChannel(&c)
		} else {
			go func(c chan struct{}, duration time.Duration) {
				ticker := time.NewTicker(duration)
				for range ticker.C {
					c <- struct{}{}
				}
			}(c, cfg.StoreInterval)
		}

		go func(ch chan struct{}) {
			for range ch {
				backup.UpdateBackup(cfg)
			}
		}(c)
	}

	handler := handlers.New(datastore)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handler.GetAllMetrics)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{typeMetric}/{nameMetric}/{value}", handler.UpdateMetric)
		r.Post("/", handler.PostUpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{typeMetric}/{nameMetric}", handler.GetMetric)
		r.Post("/", handler.PostGetMetric)
	})

	log.Fatal(http.ListenAndServe(cfg.Address, r))
}
