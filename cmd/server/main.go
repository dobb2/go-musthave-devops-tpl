package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/dobb2/go-musthave-devops-tpl/internal/logging"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"

	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/dobb2/go-musthave-devops-tpl/internal/backup"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/handlers"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := logging.CreateLogger()
	logger.Info().Msg("Start server")

	cfg := config.CreateServerConfig(logger)
	r := chi.NewRouter()

	var datastore storage.MetricGetterCreatorUpdater
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		logger.Warn().Err(err).Msg("cannot connect to database")
	}
	defer db.Close()

	if cfg.DatabaseDSN != "" {
		datastore, err = postgres.Create(db)
		if err != nil {
			logger.Fatal().Err(err).Msg("cannot create table in database")
		}
	} else {
		datastore = cache.Create()
	}

	if cfg.DatabaseDSN == "" && cfg.StoreFile != "" {
		backup := backup.New(datastore)
		if cfg.Restore {
			err = backup.Restore(cfg)
			if err != nil {
				logger.Warn().Err(err).Msg("failed backup to cache")
			} else {
				logger.Info().Msg("backup success")
			}
		}

		c := make(chan struct{})

		if cfg.StoreInterval == 0 {
			datastore.AddChannel(&c)
			logger.Info().Msg("add channel")
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
				err = backup.UpdateBackup(cfg)
				if err != nil {
					logger.Warn().Err(err).Msg("failed update backup file")
				}
			}
		}(c)
	}

	handler := handlers.New(datastore, logger)

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

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", handler.PostUpdateBatchMetrics)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{typeMetric}/{nameMetric}", handler.GetMetric)
		r.Post("/", handler.PostGetMetric)
	})

	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		logger.Fatal().Err(err).Msg("server failed")
	}
}
