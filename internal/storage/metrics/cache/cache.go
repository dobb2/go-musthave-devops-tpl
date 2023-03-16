package cache

import (
	"errors"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type Metrics struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}

func Create() Metrics {
	return Metrics{
		GaugeMetrics:   map[string]float64{},
		CounterMetrics: map[string]int64{},
	}
}

func (m Metrics) UpdateGauge(typeMetric string, value float64) {
	m.GaugeMetrics[typeMetric] = value
}

func (m Metrics) UpdateCounter(typeMetric string, value int64) {
	m.CounterMetrics[typeMetric] = value
}

func (m Metrics) GetAllMetrics() ([]metrics.Metric, error) {
	countMetrics := len(m.GaugeMetrics) + len(m.CounterMetrics)
	c := make([]metrics.Metric, 0, countMetrics)

	if countMetrics == 0 {
		return c, errors.New("No metrics")
	}

	for NameMetric, ValueMetric := range m.GaugeMetrics { // Порядок не определен
		metric := metrics.Metric{
			TypeMetric: "gauge",
			NameMetric: NameMetric,
			Value:      strconv.FormatFloat(ValueMetric, 'f', -1, 64),
		}
		c = append(c, metric)
	}
	for NameMetric, ValueMetric := range m.CounterMetrics { // Порядок не определен
		metric := metrics.Metric{
			TypeMetric: "gauge",
			NameMetric: NameMetric,
			Value:      strconv.FormatInt(ValueMetric, 10),
		}
		c = append(c, metric)
	}
	return c, nil
}

func (m Metrics) GetValue(typeMetric string, NameMetric string) (metrics.Metric, error) {
	switch typeMetric {
	case "gauge":
		if value, ok := m.GaugeMetrics[NameMetric]; ok == true {
			ResultMetric := metrics.Metric{
				TypeMetric: typeMetric,
				NameMetric: NameMetric,
				Value:      strconv.FormatFloat(value, 'f', -1, 64),
			}
			return ResultMetric, nil
		}
		return metrics.Metric{}, errors.New("unknown metric")
	case "counter":
		if value, ok := m.CounterMetrics[NameMetric]; ok == true {
			ResultMetric := metrics.Metric{
				TypeMetric: typeMetric,
				NameMetric: NameMetric,
				Value:      strconv.FormatInt(value, 10),
			}
			return ResultMetric, nil
		}
		return metrics.Metric{}, errors.New("unknown metric")
	default:
		return metrics.Metric{}, errors.New("Invalid type metric")
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
	m.CounterMetrics["PollCount"]++
	m.UpdateGauge("RandomValue", float64(rand.Float64()))
}
