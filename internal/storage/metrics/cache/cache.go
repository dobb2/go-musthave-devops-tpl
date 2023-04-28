package cache

import (
	"errors"
	"math/rand"
	"runtime"
	"time"

	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
)

type Metrics struct {
	Metrics map[string]metrics.Metrics `json:"metrics"`
	сhan    *chan struct{}
}

func Create() *Metrics {
	return &Metrics{
		Metrics: map[string]metrics.Metrics{},
		сhan:    nil,
	}
}

func (m Metrics) UpdateMetrics(metrics []metrics.Metrics) error {
	for _, metric := range metrics {
		m.Metrics[metric.ID] = metric
	}
	return nil
}

func (m *Metrics) AddChannel(c *chan struct{}) error {
	m.сhan = c
	return nil
}

func (m Metrics) UpdateGauge(nameMetric string, value float64) error {
	Value := value
	metric := metrics.Metrics{
		ID:    nameMetric,
		MType: "gauge",
		Value: &Value,
	}
	m.Metrics[nameMetric] = metric
	if m.сhan != nil {
		*m.сhan <- struct{}{}
	}

	return nil
}

func (m Metrics) UpdateCounter(nameMetric string, value int64) error {
	Delta := value
	_, ok := m.Metrics[nameMetric]
	if ok {
		Delta += *m.Metrics[nameMetric].Delta
	}
	metric := metrics.Metrics{
		ID:    nameMetric,
		MType: "counter",
		Delta: &Delta,
	}
	m.Metrics[nameMetric] = metric
	if m.сhan != nil {
		*m.сhan <- struct{}{}
	}

	return nil
}

func (m Metrics) GetAllMetrics() ([]metrics.Metrics, error) {
	countMetrics := len(m.Metrics)
	allMetric := make([]metrics.Metrics, 0, countMetrics)

	if countMetrics == 0 {
		return allMetric, errors.New("no metrics")
	}

	for _, Metric := range m.Metrics {
		allMetric = append(allMetric, Metric)
	}

	return allMetric, nil
}

func (m Metrics) GetValue(typeMetric string, NameMetric string) (metrics.Metrics, error) {
	switch typeMetric {
	case "gauge":
		if metric, ok := m.Metrics[NameMetric]; ok {
			Value := *metric.Value
			ResultMetric := metrics.Metrics{
				ID:    NameMetric,
				MType: typeMetric,
				Value: &Value,
			}
			return ResultMetric, nil
		}
		return metrics.Metrics{}, errors.New("unknown metric")
	case "counter":
		if metric, ok := m.Metrics[NameMetric]; ok {
			Delta := *metric.Delta
			ResultMetric := metrics.Metrics{
				ID:    NameMetric,
				MType: typeMetric,
				Delta: &Delta,
			}
			return ResultMetric, nil
		}
		return metrics.Metrics{}, errors.New("unknown metric")
	default:
		return metrics.Metrics{}, errors.New("invalid type metric")
	}
}

func (m Metrics) CollectMetrics() {
	rand.Seed(time.Now().UnixNano())
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	// Misc memory stats
	m.UpdateGauge("Alloc", float64(rtm.Alloc))
	m.UpdateGauge("Sys", float64(rtm.Sys))
	m.UpdateGauge("TotalAlloc", float64(rtm.TotalAlloc))
	m.UpdateGauge("Mallocs", float64(rtm.Mallocs))
	m.UpdateGauge("Frees", float64(rtm.Frees))
	m.UpdateGauge("Alloc", float64(rtm.Alloc))
	m.UpdateGauge("BuckHashSys", float64(rtm.BuckHashSys))
	m.UpdateGauge("Frees", float64(rtm.Frees))
	m.UpdateGauge("GCCPUFraction", float64(rtm.GCCPUFraction))
	m.UpdateGauge("GCSys", float64(rtm.GCSys))
	m.UpdateGauge("HeapAlloc", float64(rtm.HeapAlloc))
	m.UpdateGauge("HeapIdle", float64(rtm.HeapIdle))
	m.UpdateGauge("HeapInuse", float64(rtm.HeapInuse))
	m.UpdateGauge("HeapObjects", float64(rtm.HeapObjects))
	m.UpdateGauge("HeapReleased", float64(rtm.HeapReleased))
	m.UpdateGauge("HeapSys", float64(rtm.HeapSys))
	m.UpdateGauge("LastGC", float64(rtm.LastGC))
	m.UpdateGauge("Lookups", float64(rtm.Lookups))
	m.UpdateGauge("MCacheInuse", float64(rtm.MCacheInuse))
	m.UpdateGauge("MCacheSys", float64(rtm.MCacheSys))
	m.UpdateGauge("MSpanInuse", float64(rtm.MSpanInuse))
	m.UpdateGauge("MSpanSys", float64(rtm.MSpanSys))
	m.UpdateGauge("Mallocs", float64(rtm.Mallocs))
	m.UpdateGauge("NextGC", float64(rtm.NextGC))
	m.UpdateGauge("NumForcedGC", float64(rtm.NumForcedGC))
	m.UpdateGauge("NumGC", float64(rtm.NumGC))
	m.UpdateGauge("OtherSys", float64(rtm.OtherSys))
	m.UpdateGauge("PauseTotalNs", float64(rtm.PauseTotalNs))
	m.UpdateGauge("StackInuse", float64(rtm.StackInuse))
	m.UpdateGauge("StackSys", float64(rtm.StackSys))
	m.UpdateGauge("Sys", float64(rtm.Sys))
	m.UpdateGauge("TotalAlloc", float64(rtm.TotalAlloc))
	m.UpdateCounter("PollCount", 1)
	m.UpdateGauge("RandomValue", float64(rand.Float64()))
}

func (m Metrics) GetPing() error {
	return nil
}
