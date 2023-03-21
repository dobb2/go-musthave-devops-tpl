package storage

import "github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"

type MetricCreatorUpdater interface {
	UpdateGauge(typeMetric string, value float64)
	UpdateCounter(typeMetric string, value int64)
	GetValue(typeMetric string, NameMetric string) (metrics.Metric, error)
	GetAllMetrics() ([]metrics.Metric, error)
}
