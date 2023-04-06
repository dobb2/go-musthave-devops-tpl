package storage

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
)

type MetricCreatorUpdater interface {
	UpdateGauge(typeMetric string, value float64)
	UpdateCounter(typeMetric string, value int64)
	GetValue(typeMetric string, NameMetric string) (metrics.Metrics, error)
	GetAllMetrics() ([]metrics.Metrics, error)
	GetPing() error
	UploadMetrics([]metrics.Metrics)
	AddChannel(*chan struct{})
}
