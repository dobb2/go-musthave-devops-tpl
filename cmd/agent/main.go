package main

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/logging"
	"time"

	"github.com/dobb2/go-musthave-devops-tpl/internal/client"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
)

func main() {
	logger := logging.CreateLogger()
	logger.Info().Msg("Start agent")
	cfg := config.CreateAgentConfig(logger)
	m := cache.Create()

	tickerCollector := time.NewTicker(cfg.PollInterval)
	tickerSender := time.NewTicker(cfg.ReportInterval)

	for {
		select {
		case <-tickerCollector.C:
			m.CollectMetrics()
		case <-tickerSender.C:
			client.PutMetric(m, cfg, logger)
		}
	}
}
