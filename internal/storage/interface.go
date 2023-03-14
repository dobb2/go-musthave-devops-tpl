package storage

type MetricCreatorUpdater interface {
	UpdateGauge(typeMetric string, value float64)
	UpdateCounter(typeMetric string, value int64)
}
