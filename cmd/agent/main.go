package main

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/client"
	"github.com/dobb2/go-musthave-devops-tpl/internal/logging"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"github.com/shirou/gopsutil/v3/mem"
	"time"

	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
)

func main() {
	v, _ := mem.VirtualMemory()
	logger := logging.CreateLogger()
	logger.Info().Msg("Start agent")
	cfg := config.CreateAgentConfig(logger)
	logger.Info().Msg("Key" + cfg.Key)
	m := cache.Create()

	agent := client.New(m, logger, cfg)

	tickerCollector := time.NewTicker(cfg.PollInterval)
	tickerSender := time.NewTicker(cfg.ReportInterval)

	MetricCh := make(chan metrics.Metrics)
	defer close(MetricCh)

	go agent.WorkPool(MetricCh)
	for {
		select {
		case <-tickerCollector.C:
			go agent.CollectMetrics(v)
		case <-tickerSender.C:
			go agent.PutMetric(MetricCh)
		}
	}

}
