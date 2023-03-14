package storage

type Metrics struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		GaugeMetrics:   map[string]float64{},
		CounterMetrics: map[string]int64{},
	}
}
