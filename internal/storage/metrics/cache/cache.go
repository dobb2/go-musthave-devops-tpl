package cache

import (
	"errors"
	"sync"

	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
)

type Metrics struct {
	mu      sync.RWMutex
	Metrics map[string]metrics.Metrics `json:"metrics"`
	сhan    *chan struct{}
}

func Create() *Metrics {
	return &Metrics{
		Metrics: map[string]metrics.Metrics{},
		сhan:    nil,
	}
}

func (m *Metrics) UpdateMetrics(metrics []metrics.Metrics) error {

	for _, metric := range metrics {
		if metric.MType == "counter" {
			m.UpdateCounter(metric.ID, *metric.Delta)
		} else if metric.MType == "gauge" {
			m.UpdateGauge(metric.ID, *metric.Value)
		} else {
			return errors.New("unknown metric type")
		}
	}
	return nil
}

func (m *Metrics) AddChannel(c *chan struct{}) error {
	m.сhan = c
	return nil
}

func (m *Metrics) UpdateGauge(nameMetric string, value float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
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

func (m *Metrics) UpdateCounter(nameMetric string, value int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

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

func (m *Metrics) GetAllMetrics() ([]metrics.Metrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

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

func (m *Metrics) GetValue(typeMetric string, NameMetric string) (metrics.Metrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

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

func (m *Metrics) GetPing() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return nil
}
