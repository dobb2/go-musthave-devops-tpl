package postgres

import (
	"database/sql"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"log"
)

type MetricsStorer struct {
	db *sql.DB
}

func New(db *sql.DB) MetricsStorer {
	return MetricsStorer{db: db}
}

func (m MetricsStorer) UpdateGauge(typeMetric string, value float64) {
	log.Println("test gauge")
}
func (m MetricsStorer) UpdateCounter(typeMetric string, value int64) {
	log.Println("test counter")
}

func (m MetricsStorer) GetValue(typeMetric string, NameMetric string) (metrics.Metrics, error) {
	return metrics.Metrics{}, nil
}

func (m MetricsStorer) GetAllMetrics() ([]metrics.Metrics, error) {
	return make([]metrics.Metrics, 0), nil
}

func (m MetricsStorer) GetPing() error {
	return m.db.Ping()
}
