package main

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/client"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"time"
)

func main() {
	cfg := config.CreateAgentConfig()

	m := cache.Create()

	ticker := time.NewTicker(cfg.PollInterval)
	ticker2 := time.NewTicker(cfg.ReportInterval)

	for {
		select {
		case <-ticker.C:
			m.CollectMetrics()
		case <-ticker2.C:
			client.PutMetric(&m, cfg)
		}
	}
}
