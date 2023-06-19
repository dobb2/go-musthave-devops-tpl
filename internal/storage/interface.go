package storage

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
)

type MetricGetterCreatorUpdater interface {
	UpdateGauge(string, float64) error
	UpdateCounter(string, int64) error
	GetValue(string, string) (metrics.Metrics, error)
	GetAllMetrics() ([]metrics.Metrics, error)
	GetPing() error
	UpdateMetrics([]metrics.Metrics) error
	AddChannel(*chan struct{}) error
}
