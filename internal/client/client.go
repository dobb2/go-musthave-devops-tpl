package client

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"runtime"
	"time"

	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/crypto"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"github.com/go-resty/resty/v2"
)

type MetricCreatorUpdater interface {
	UpdateGauge(string, float64) error
	UpdateCounter(string, int64) error
	GetAllMetrics() ([]metrics.Metrics, error)
}

type MetricsАgent struct {
	cache  MetricCreatorUpdater
	logger zerolog.Logger
	config config.AgentConfig
}

func New(metrics MetricCreatorUpdater, logger zerolog.Logger, config config.AgentConfig) MetricsАgent {
	return MetricsАgent{
		cache:  metrics,
		logger: logger,
		config: config,
	}
}

func (m MetricsАgent) SendBatchMetric(metrics []metrics.Metrics) {
	client := resty.New().
		SetBaseURL("http://" + m.config.Address).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	out, err := json.Marshal(metrics)
	if err != nil {
		m.logger.Error().Err(err).Msg("unsuccessful marshal metrics to json")
		return
	}
	m.logger.Info().Msg("send metric")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(out).
		Post("/updates/")

	if err != nil {
		m.logger.Error().Err(err).Msg("unsuccessful request")
	} else {
		m.logger.Info().Msgf("Status code: %d", resp.StatusCode())
	}
}

func (m MetricsАgent) SendMetric(metric metrics.Metrics) {
	client := resty.New().
		SetBaseURL("http://" + m.config.Address).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	out, err := json.Marshal(metric)
	if err != nil {
		m.logger.Error().Err(err).Msg("unsuccessful marshal metrics to json")
		return
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(out).
		Post("/update/")

	if err != nil {
		m.logger.Error().Err(err).Msg("unsuccessful request")
	} else {
		m.logger.Info().Msgf("Status code: %d", resp.StatusCode())
	}
}

func (m *MetricsАgent) PutMetric(inputCh chan<- metrics.Metrics) {
	cacheMetrics, _ := m.cache.GetAllMetrics()
	for _, Metric := range cacheMetrics { // Порядок не определен
		switch Metric.MType {
		case "counter":
			Metric.Hash = crypto.Hash(fmt.Sprintf("%s:counter:%d", Metric.ID, *Metric.Delta), m.config.Key)
		case "gauge":
			Metric.Hash = crypto.Hash(fmt.Sprintf("%s:gauge:%f", Metric.ID, *Metric.Value), m.config.Key)
		default:
			m.logger.Warn().Msg("invalid type metric for create hash")
		}
		inputCh <- Metric
	}
}

func (m *MetricsАgent) WorkPool(inputCh <-chan metrics.Metrics) {
	if m.config.RateLimit == 0 {
		buf := make([]metrics.Metrics, 0, m.config.MetricMaxAmount)
		for metric := range inputCh {
			buf = append(buf, metric)
			if len(buf) == m.config.MetricMaxAmount {
				m.logger.Info().Msg("tuta mi kak voobshe metric")
				m.SendBatchMetric(buf)
				buf = buf[:0]
			}
		}
	} else {
		workersCount := m.config.RateLimit
		for i := 0; i < workersCount; i++ {
			go func() {
				for metric := range inputCh {
					m.SendMetric(metric)
				}
			}()
		}
	}
}

func (m MetricsАgent) CollectMetrics(v *mem.VirtualMemoryStat) {
	rand.Seed(time.Now().UnixNano())
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	// Misc memory stats
	m.cache.UpdateGauge("Alloc", float64(rtm.Alloc))
	m.cache.UpdateGauge("Sys", float64(rtm.Sys))
	m.cache.UpdateGauge("TotalAlloc", float64(rtm.TotalAlloc))
	m.cache.UpdateGauge("Mallocs", float64(rtm.Mallocs))
	m.cache.UpdateGauge("Frees", float64(rtm.Frees))
	m.cache.UpdateGauge("Alloc", float64(rtm.Alloc))
	m.cache.UpdateGauge("BuckHashSys", float64(rtm.BuckHashSys))
	m.cache.UpdateGauge("Frees", float64(rtm.Frees))
	m.cache.UpdateGauge("GCCPUFraction", float64(rtm.GCCPUFraction))
	m.cache.UpdateGauge("GCSys", float64(rtm.GCSys))
	m.cache.UpdateGauge("HeapAlloc", float64(rtm.HeapAlloc))
	m.cache.UpdateGauge("HeapIdle", float64(rtm.HeapIdle))
	m.cache.UpdateGauge("HeapInuse", float64(rtm.HeapInuse))
	m.cache.UpdateGauge("HeapObjects", float64(rtm.HeapObjects))
	m.cache.UpdateGauge("HeapReleased", float64(rtm.HeapReleased))
	m.cache.UpdateGauge("HeapSys", float64(rtm.HeapSys))
	m.cache.UpdateGauge("LastGC", float64(rtm.LastGC))
	m.cache.UpdateGauge("Lookups", float64(rtm.Lookups))
	m.cache.UpdateGauge("MCacheInuse", float64(rtm.MCacheInuse))
	m.cache.UpdateGauge("MCacheSys", float64(rtm.MCacheSys))
	m.cache.UpdateGauge("MSpanInuse", float64(rtm.MSpanInuse))
	m.cache.UpdateGauge("MSpanSys", float64(rtm.MSpanSys))
	m.cache.UpdateGauge("Mallocs", float64(rtm.Mallocs))
	m.cache.UpdateGauge("NextGC", float64(rtm.NextGC))
	m.cache.UpdateGauge("NumForcedGC", float64(rtm.NumForcedGC))
	m.cache.UpdateGauge("NumGC", float64(rtm.NumGC))
	m.cache.UpdateGauge("OtherSys", float64(rtm.OtherSys))
	m.cache.UpdateGauge("PauseTotalNs", float64(rtm.PauseTotalNs))
	m.cache.UpdateGauge("StackInuse", float64(rtm.StackInuse))
	m.cache.UpdateGauge("StackSys", float64(rtm.StackSys))
	m.cache.UpdateGauge("Sys", float64(rtm.Sys))
	m.cache.UpdateGauge("TotalAlloc", float64(rtm.TotalAlloc))
	m.cache.UpdateCounter("PollCount", 1)
	m.cache.UpdateGauge("RandomValue", float64(rand.Float64()))
	m.cache.UpdateGauge("TotalMemory", float64(v.Total))
	m.cache.UpdateGauge("FreeMemory", float64(v.Free))
	m.cache.UpdateGauge("CPUutilization1", float64(runtime.NumCPU()))

}
