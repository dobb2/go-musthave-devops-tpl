package main

import (
	"database/sql"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/database"
	_ "github.com/jackc/pgx"
	"log"
	"net/http"
	"time"

	"github.com/dobb2/go-musthave-devops-tpl/internal/backup"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/handlers"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.CreateServerConfig()
	r := chi.NewRouter()
	var datastore storage.MetricCreatorUpdater

	if cfg.DatabaseDSN != "" {
		db, err := sql.Open("pgx",
			cfg.DatabaseDSN)
		if err != nil {
			log.Println(err)
		}
		defer db.Close()
		datastore = database.New(db)
	} else {
		datastore = cache.Create()
	}

	if cfg.StoreFile != "" {
		c := make(chan struct{})
		backup := backup.New(datastore)
		if cfg.Restore {
			backup.Restore(cfg)
		}

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
	r.Use(middleware.WithValue("Key", cfg.Key))
	r.Use(middleware.Compress(5))

	r.Get("/", handler.GetAllMetrics)
	r.Get("/ping", handler.GetPing)

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
